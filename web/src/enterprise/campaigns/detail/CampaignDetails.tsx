import slugify from 'slugify'
import { LoadingSpinner } from '@sourcegraph/react-loading-spinner'
import AlertCircleIcon from 'mdi-react/AlertCircleIcon'
import React, { useState, useEffect, useRef, useMemo, useCallback } from 'react'
import * as GQL from '../../../../../shared/src/graphql/schema'
import { HeroPage } from '../../../components/HeroPage'
import { PageTitle } from '../../../components/PageTitle'
import { UserAvatar } from '../../../user/UserAvatar'
import { Timestamp } from '../../../components/time/Timestamp'
import { noop, isEqual } from 'lodash'
import { Form } from '../../../components/Form'
import {
    fetchCampaignById,
    updateCampaign,
    deleteCampaign,
    createCampaign,
    closeCampaign,
    fetchPatchSetById,
    queryPatchesFromCampaign,
    queryPatchesFromPatchSet,
    queryChangesets,
    queryPatchFileDiffs,
} from './backend'
import { useError, useObservable } from '../../../../../shared/src/util/useObservable'
import { asError } from '../../../../../shared/src/util/errors'
import * as H from 'history'
import { CampaignBurndownChart } from './BurndownChart'
import { AddChangesetForm } from './AddChangesetForm'
import { Subject, of, merge, Observable, NEVER } from 'rxjs'
import { renderMarkdown, highlightCodeSafe } from '../../../../../shared/src/util/markdown'
import { ErrorAlert } from '../../../components/alerts'
import { Markdown } from '../../../../../shared/src/components/Markdown'
import { switchMap, distinctUntilChanged } from 'rxjs/operators'
import { ThemeProps } from '../../../../../shared/src/theme'
import { CampaignDescriptionField } from './form/CampaignDescriptionField'
import { CampaignStatus } from './CampaignStatus'
import { CampaignUpdateDiff } from './CampaignUpdateDiff'
import { CampaignActionsBar } from './CampaignActionsBar'
import { CampaignTitleField } from './form/CampaignTitleField'
import { CampaignChangesets } from './changesets/CampaignChangesets'
import { CampaignDiffStat } from './CampaignDiffStat'
import { pluralize } from '../../../../../shared/src/util/strings'
import { ExtensionsControllerProps } from '../../../../../shared/src/extensions/controller'
import { PlatformContextProps } from '../../../../../shared/src/platform/context'
import { TelemetryProps } from '../../../../../shared/src/telemetry/telemetryService'
import { CampaignPatches } from './patches/CampaignPatches'
import { PatchSetPatches } from './patches/PatchSetPatches'
import { CampaignBranchField } from './form/CampaignBranchField'
import { repeatUntil } from '../../../../../shared/src/util/rxjs/repeatUntil'
import { MinimalCampaign, MinimalPatchSet } from './CampaignArea'

export type CampaignUIMode = 'viewing' | 'deleting' | 'closing'

interface Props extends ThemeProps, ExtensionsControllerProps, PlatformContextProps, TelemetryProps {
    campaign: MinimalCampaign
    authenticatedUser: Pick<GQL.IUser, 'id' | 'username' | 'avatarURL'>
    history: H.History
    location: H.Location

    /** For testing only. */
    _fetchPatchSetById?: typeof fetchPatchSetById | ((patchSet: GQL.ID) => Observable<MinimalPatchSet | null>)
    /** For testing only. */
    _queryPatchesFromCampaign?: typeof queryPatchesFromCampaign
    /** For testing only. */
    _queryPatchesFromPatchSet?: typeof queryPatchesFromPatchSet
    /** For testing only. */
    _queryPatchFileDiffs?: typeof queryPatchFileDiffs
    /** For testing only. */
    _queryChangesets?: typeof queryChangesets
    /** For testing only. */
    _noSubject?: boolean
}

/**
 * The area for a single campaign.
 */
export const CampaignDetails: React.FunctionComponent<Props> = ({
    campaign,
    history,
    location,
    authenticatedUser,
    isLightTheme,
    extensionsController,
    platformContext,
    telemetryService,
    _fetchPatchSetById = fetchPatchSetById,
    _queryPatchesFromCampaign = queryPatchesFromCampaign,
    _queryPatchesFromPatchSet = queryPatchesFromPatchSet,
    _queryPatchFileDiffs = queryPatchFileDiffs,
    _queryChangesets = queryChangesets,
}) => {
    /** Retrigger campaign fetching */
    const campaignUpdates = useMemo(() => new Subject<void>(), [])
    /** Retrigger changeset fetching */
    const changesetUpdates = useMemo(() => new Subject<void>(), [])

    const [mode, setMode] = useState<CampaignUIMode>(campaignID ? 'viewing' : 'editing')

    // To report errors from saving or deleting
    const [alertError, setAlertError] = useState<Error>()

    const patchSetID: GQL.ID | null = new URLSearchParams(location.search).get('patchSet')
    useEffect(() => {
        if (patchSetID) {
            setMode('editing')
        }
    }, [patchSetID])

    const patchSet = useObservable(
        useMemo(() => (!patchSetID ? NEVER : _fetchPatchSetById(patchSetID)), [patchSetID, _fetchPatchSetById])
    )

    const onAddChangeset = useCallback((): void => {
        // we also check the campaign.changesets.totalCount, so an update to the campaign is required as well
        campaignUpdates.next()
        changesetUpdates.next()
    }, [campaignUpdates, changesetUpdates])

    // Patch set was not found
    if (patchSet === null) {
        return <HeroPage icon={AlertCircleIcon} title="Patch set not found" />
    }

    // On update, check if an update is possible
    if (!!campaign && !!patchSet) {
        if (!campaign.patchSet?.id) {
            return <HeroPage icon={AlertCircleIcon} title="Cannot update a manual campaign with a patch set" />
        }
        if (campaign.closedAt) {
            return <HeroPage icon={AlertCircleIcon} title="Cannot update a closed campaign" />
        }
    }

    const specifyingBranchAllowed =
        // on campaign creation
        (!campaign && patchSet) ||
        // or when no changesets have been published or are being published as well
        (campaign &&
            campaign.changesets.totalCount === 0 &&
            campaign.status.state !== GQL.BackgroundProcessState.PROCESSING)

    const onClose = async (closeChangesets: boolean): Promise<void> => {
        if (!confirm('Are you sure you want to close the campaign?')) {
            return
        }
        setMode('closing')
        try {
            await closeCampaign(campaign.id, closeChangesets)
            campaignUpdates.next()
        } catch (error) {
            setAlertError(asError(error))
        } finally {
            setMode('viewing')
        }
    }

    const onDelete = async (closeChangesets: boolean): Promise<void> => {
        if (!confirm('Are you sure you want to delete the campaign?')) {
            return
        }
        setMode('deleting')
        try {
            await deleteCampaign(campaign.id, closeChangesets)
            history.push('/campaigns')
        } catch (error) {
            setAlertError(asError(error))
            setMode('viewing')
        }
    }

    const afterRetry = (updatedCampaign: MinimalCampaign): void => {
        setCampaign(updatedCampaign)
        campaignUpdates.next()
    }

    const author = campaign ? campaign.author : authenticatedUser

    const totalChangesetCount = campaign?.changesets.totalCount ?? 0

    const totalPatchCount = (campaign?.patches.totalCount ?? 0) + (patchSet?.patches.totalCount ?? 0)

    const campaignFormID = 'campaign-form'

    return (
        <>
            <PageTitle title={campaign ? campaign.name : 'New campaign'} />
            <CampaignActionsBar
                previewingPatchSet={!!patchSet}
                mode={mode}
                campaign={campaign}
                onEdit={onEdit}
                onClose={onClose}
                onDelete={onDelete}
                formID={campaignFormID}
            />
            {alertError && <ErrorAlert error={alertError} history={history} />}
            {campaign && !patchSet && !['saving', 'editing'].includes(mode) && (
                <CampaignStatus campaign={campaign} afterRetry={afterRetry} history={history} />
            )}
            <Form id={campaignFormID} onSubmit={onSubmit} onReset={onCancel} className="e2e-campaign-form">
                {['saving', 'editing'].includes(mode) && (
                    <>
                        <h3>Details</h3>
                        {/* Existing non-manual campaign, but not updating with a new set of patches */}
                        {campaign && !!campaign.patchSet && !patchSet && (
                            <div className="card">
                                <div className="card-body">
                                    <h3 className="card-title">Want to update the patches?</h3>
                                    <p>
                                        Using the{' '}
                                        <a
                                            href="https://github.com/sourcegraph/src-cli"
                                            rel="noopener noreferrer"
                                            target="_blank"
                                        >
                                            src CLI
                                        </a>
                                        , you can also apply a new patch set to an existing campaign. Following the
                                        creation of a new patch set that contains new patches, with the
                                    </p>
                                    <div className="alert alert-secondary">
                                        <code
                                            dangerouslySetInnerHTML={{
                                                __html: highlightCodeSafe(
                                                    '$ src action exec -f action.json | src campaign patchset create-from-patches',
                                                    'bash'
                                                ),
                                            }}
                                        />
                                    </div>
                                    <p>
                                        command, a URL will be output that will guide you to the web UI to allow you to
                                        change an existing campaign’s patch set.
                                    </p>
                                    <p className="mb-0">
                                        Take a look at the{' '}
                                        <a
                                            href="https://docs.sourcegraph.com/user/campaigns/updating_campaigns"
                                            rel="noopener noreferrer"
                                            target="_blank"
                                        >
                                            documentation on updating campaigns
                                        </a>{' '}
                                        for more information.
                                    </p>
                                </div>
                            </div>
                        )}
                    </>
                )}
                {/* If we are in the update view */}
                {campaign && patchSet && (
                    <>
                        <CampaignUpdateDiff
                            campaign={campaign}
                            patchSet={patchSet}
                            queryPatchFileDiffs={queryPatchFileDiffs}
                            history={history}
                            location={location}
                            isLightTheme={isLightTheme}
                            className="mt-4"
                        />
                        <div className="mb-0">
                            <button
                                type="reset"
                                form={campaignFormID}
                                className="btn btn-secondary mr-1"
                                onClick={onCancel}
                                disabled={mode !== 'editing'}
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                form={campaignFormID}
                                className="btn btn-primary"
                                disabled={mode !== 'editing' || patchSet?.patches.totalCount === 0}
                            >
                                Update
                            </button>
                        </div>
                    </>
                )}
            </Form>

            {/* Iff either campaign XOR patchset are present */}
            {!(campaign && patchSet) && (campaign || patchSet) && (
                <>
                    {campaign && !['saving', 'editing'].includes(mode) && (
                        <>
                            <div className="card mt-2">
                                <div className="card-header">
                                    <strong>
                                        <UserAvatar user={author} className="icon-inline" /> {author.username}
                                    </strong>{' '}
                                    started <Timestamp date={campaign.createdAt} />
                                </div>
                                <div className="card-body">
                                    <Markdown
                                        dangerousInnerHTML={renderMarkdown(campaign.description || '_No description_')}
                                        history={history}
                                    />
                                </div>
                            </div>
                            {totalChangesetCount > 0 && (
                                <>
                                    <h3 className="mt-4 mb-2">Progress</h3>
                                    <CampaignBurndownChart
                                        changesetCountsOverTime={campaign.changesetCountsOverTime}
                                        history={history}
                                    />
                                </>
                            )}

                            {/* Only campaigns that have no patch set can add changesets manually. */}
                            {!campaign.patchSet && campaign.viewerCanAdminister && !campaign.closedAt && (
                                <>
                                    {totalChangesetCount === 0 && (
                                        <div className="mt-4 mb-2 alert alert-info e2e-campaign-get-started">
                                            Add a changeset to get started.
                                        </div>
                                    )}
                                    <AddChangesetForm
                                        campaignID={campaign.id}
                                        onAdd={onAddChangeset}
                                        history={history}
                                    />
                                </>
                            )}
                        </>
                    )}

                    {totalChangesetCount + totalPatchCount > 0 && (
                        <>
                            <h3 className="mt-4 d-flex align-items-end mb-0">
                                {totalPatchCount > 0 && (
                                    <>
                                        {totalPatchCount} {pluralize('Patch', totalPatchCount, 'Patches')}
                                    </>
                                )}
                                {(totalChangesetCount > 0 || !!campaign) && totalPatchCount > 0 && (
                                    <span className="mx-1">/</span>
                                )}
                                {(totalChangesetCount > 0 || !!campaign) && (
                                    <>
                                        {totalChangesetCount} {pluralize('Changeset', totalChangesetCount)}
                                    </>
                                )}{' '}
                                {(patchSet || campaign) && (
                                    <CampaignDiffStat campaign={campaign} patchSet={patchSet} className="ml-2 mb-0" />
                                )}
                            </h3>
                            {totalPatchCount > 0 &&
                                (campaign ? (
                                    <CampaignPatches
                                        campaign={campaign}
                                        campaignUpdates={campaignUpdates}
                                        changesetUpdates={changesetUpdates}
                                        enablePublishing={!campaign.closedAt}
                                        queryPatchesFromCampaign={_queryPatchesFromCampaign}
                                        queryPatchFileDiffs={_queryPatchFileDiffs}
                                        history={history}
                                        location={location}
                                        isLightTheme={isLightTheme}
                                    />
                                ) : (
                                    <PatchSetPatches
                                        patchSet={patchSet!}
                                        campaignUpdates={campaignUpdates}
                                        changesetUpdates={changesetUpdates}
                                        enablePublishing={false}
                                        queryPatchesFromPatchSet={_queryPatchesFromPatchSet}
                                        queryPatchFileDiffs={_queryPatchFileDiffs}
                                        history={history}
                                        location={location}
                                        isLightTheme={isLightTheme}
                                    />
                                ))}
                            {totalChangesetCount > 0 && (
                                <CampaignChangesets
                                    campaign={campaign}
                                    changesetUpdates={changesetUpdates}
                                    campaignUpdates={campaignUpdates}
                                    queryChangesets={_queryChangesets}
                                    history={history}
                                    location={location}
                                    isLightTheme={isLightTheme}
                                    extensionsController={extensionsController}
                                    platformContext={platformContext}
                                    telemetryService={telemetryService}
                                />
                            )}
                        </>
                    )}
                </>
            )}
        </>
    )
}
