package db_helpers

import (
	"errors"
	"time"
)

func validateBorrow(book Book) error {
	if book.Student.Name == "" {
		return errors.New("student details not provided")
	}
	books := Find(book.Name)
	if len(books) == 0 {
		return errors.New("book not found")
	} else if len(books) > 1 {
		return errors.New("given name matched to more than one book")
	} else {
		if !(books[0].Available) {
			return errors.New("book is not available")
		}
	}
	return nil
}

func validateReturn(book Book) error {
	books := Find(book.Name)
	if len(books) == 0 {
		return errors.New("book not found")
	} else if len(books) > 1 {
		return errors.New("given name matched to more than one book")
	} else {
		if books[0].Available {
			return errors.New("book is already available and doesn't need to be returned")
		}
	}
	return nil
}

func isReturnDelayed(book Book, days int) (int, error) {
	book = Find(book.Name)[0]
	returnObj, err := time.Parse("2006-01-02", book.ReturnDate)
	if err != nil { return 0, err }

	diff := int(time.Now().Local().Sub(returnObj).Hours()/24)
	if diff > 0 { return diff, nil }

	return 0, nil
}
