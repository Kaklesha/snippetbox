package models

func (m *SnippetModel) ExampleTransaction() error {
	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in-progress database transaction.
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// 	stmt := `INSERT INTO snippets (title, content, created, expires)
	// VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Defer a call to tx.Rollback() to ensure it is always called before the
	// function returns. If the transaction succeeds it will be already be
	// committed by the time tx.Rollback() is called, making tx.Rollback()
	// no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
	// the changes before the function returns.
	defer tx.Rollback()
	// Call Exec() on the transaction, passing in your statement and any
	// parameters. It's important to notice that tx.Exec() is called on the
	// transaction object just created, NOT the connection pool. Although we're
	// using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
	// exactly the same way.

	//Create some valiables holding dummy data. We'll remove these later on
	//during the build.
	title := "LIER snail"
	content := "LIER snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 8

	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//+++++++++++++++++++++++++++++++++++++++++
	result, err := tx.Exec(stmt, title, content, expires)
	if err != nil {

		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.

	stmtTwo := `UPDATE snippets SET title = ? WHERE id = ?`

	// Саггу out another transaction in exactly the same way.
	_, err = tx.Exec(stmtTwo, "NOT_LIER", id)
	if err != nil {
		return err
	}
	// If there are no errors, the statements in the transaction can be committed
	// to the database with the tx.Commit() method.
	err = tx.Commit()
	return err
}
