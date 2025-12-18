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

``` ```