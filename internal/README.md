# **INTERNAL/MODELS** section info 
## Description

The INTERNAL directory will contain the ancillary non-application-specific code used in the project. We'll use it to hold potentially reusable code like validation helpers and the SQL database models for the project.

### Converting Types MySQL to Go
Behind the scenes of rows.Scan() your driver will automatically convert the raw output
from the SQL database to the required native Go types. So long as you’re sensible with the
types that you’re mapping between SQL and Go, these conversions should generally Just
Work. Usually:

- ```CHAR```, ```VARCHAR``` and ```TEXT``` map to ```string```.
- ```BOOLEAN``` maps to ```bool```.
- ```INT``` maps to ```int```; ```BIGINT``` maps to ```int64```.
- ```DECIMAL``` and ```NUMERIC``` map to ```float```.
- ```TIME```, ```DATE``` and ```TIMESTAMP``` map to ```time.Time```.

Note: A quirk of our MySQL driver is that we need to use the ```parseTime=true```
parameter in our DSN to force it to convert ```TIME``` and ```DATE``` fields to ```time.Time```.
Otherwise it returns these as ```[]byte``` objects. This is one of the many *driver-specific
parameters* that it offers

### Managing null values
One thing that Go doesn’t do very well is managing ```NULL``` values in database records.
Let’s pretend that the ```title``` column in our ```snippets``` table contains a ```NULL``` value in a
particular row. If we queried that row, then ```rows.Scan()``` would return the following error
because it can’t convert ```NULL``` into a string:


``` sql: Scan error on column index 1: unsupported Scan, storing driver.Value type
&lt;nil&gt; into type *string
```

Very roughly, the fix for this is to change the field that you’re scanning into from a ```string``` to
a ```sql.NullString``` type. See this gist [this gist](https://gist.github.com/alexedwards/dc3145c8e2e6d2fd6cd9) for a working example.
But, as a rule, the easiest thing to do is simply avoid ```NULL``` values altogether. Set ```NOT NULL```
constraints on all your database columns, like we have done in this book, along with
sensible ```DEFAULT``` values as necessary.

Quirk of solution it use COALESCE into SQL DB

```SELECT id, COALESCE(title, '') FROM snippets;```


One probably of solutions it is using GENERIC a example below: 
```
package main

import (
	"database/sql/driver"
	"fmt"
)

// Null — это универсальная структура для обработки NULL значений
type Null[T any] struct {
	V     T
	Valid bool
}

// Scan реализует интерфейс sql.Scanner для чтения из базы
func (n *Null[T]) Scan(value any) error {
	if value == nil {
		n.V, n.Valid = *new(T), false
		return nil
	}
	n.Valid = true
	// Используем утверждение типа для записи значения
	n.V = value.(T)
	return nil
}

// Value реализует интерфейс driver.Valuer для записи в базу
func (n Null[T]) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.V, nil
}
```
## Example has used code above in action

 ```
type User struct {
    ID    int
    Name  Null[string] // Вместо sql.NullString
    Age   Null[int64]  // Вместо sql.NullInt64
    Email Null[string]
}

func main() {
    // Пример сканирования из БД
    var u User
    err := db.QueryRow("SELECT name, age FROM users WHERE id = ?", 1).Scan(&u.Name, &u.Age)
    
    if u.Name.Valid {
        fmt.Println("Имя:", u.Name.V)
    } else {
        fmt.Println("Имя не указано (NULL)")
    }
}
 ```



``` ```