package store

import "database/sql"

// scanner is the common interface shared by *sql.Row and *sql.Rows.
type scanner interface {
	// Scan copies the values of the current row into the values pointed at by dest.
	Scan(dest ...interface{}) error
}

func ScanInt(scanner scanner) (value int, err error) {
	err = scanner.Scan(&value)
	return value, err
}

func ScanInts(rows *sql.Rows, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []int
	for rows.Next() {
		value, err := ScanInt(rows)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func ScanFirstInt(rows *sql.Rows, err error) (int, bool, error) {
	values, err := ScanInts(rows, err)
	if err != nil || len(values) == 0 {
		return 0, false, err
	}
	return values[0], true, nil
}

func xScanBytes(scanner scanner) (value []byte, err error) {
	err = scanner.Scan(&value)
	return value, err
}

func ScanBytes(rows *sql.Rows, err error) ([][]byte, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values [][]byte
	for rows.Next() {
		value, err := xScanBytes(rows)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func ScanFirstBytes(rows *sql.Rows, err error) ([]byte, bool, error) {
	values, err := ScanBytes(rows, err)
	if err != nil || len(values) == 0 {
		return nil, false, err
	}
	return values[0], true, nil
}

func ScanString(scanner scanner) (value string, err error) {
	err = scanner.Scan(&value)
	return value, err
}

func ScanStrings(rows *sql.Rows, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		value, err := ScanString(rows)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func ScanFirstString(rows *sql.Rows, err error) (string, bool, error) {
	values, err := ScanStrings(rows, err)
	if err != nil || len(values) == 0 {
		return "", false, err
	}
	return values[0], true, nil
}
