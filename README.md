# MySQL CSV to JSON

A small MySQL UDF library for encoding CSV files to JSON, with and without column headers, written in Golang.

---

## Notes!

The MySQL syntax below in the examples might look strange, with the CSV lines being on their own lines, but they are 100% equivalent to
```sql
select`csv_to_json`('Id,Name,Age\n1,Ahmad,21\n2,Ali,50');
-- and also
select`csv_to_json`('Id,Name,Age\n' '1,Ahmad,21\n' '2,Ali,50');
```
Only a single string is being passed!

---

### `csv_to_json`

Converts a CSV file string to a json string. Will return `NULL` if given `NULL`, or if given a CSV where the values per row is not the same. The first row will be used as the column headers for the JSON objects. Duplicate column header values will have numbers appended to them, starting from 2.

```sql
`csv_to_json` ( `string` )
```

 - `` `string` ``
   - The CSV file string to be encoded.

## Examples

```sql
select`csv_to_json`(
'Id,Name,Age\n'
'1,Ahmad,21\n'
'2,Ali,50');

--'[
--     {
--         "Age": "21",
--         "Id": " 1",
--         "Name": "Ahmad"
--     },
--     {
--         "Age": "50",
--         "Id": " 2",
--         "Name": "Ali"
--     }
-- ]'

select`csv_to_json`(
'Id,Name,Name\n' -- Notice the double "Name" column
'1,Ahmad,21\n'
'2,Ali,50');

--'[
--     {
--         "Id": " 1",
--         "Name": "Ahmad",
--         "Name 2": "21" <-- The second "Name" got renamed to "Name 2"
--     },
--     {
--         "Id": " 2",
--         "Name": "Ali",
--         "Name 2": "50"
--     }
-- ]'

select`csv_to_json`(
'Id,Name,Age,Extra\n' -- Notice this row has more values than the others
'1,Ahmad,21\n'
'2,Ali,50');             -- NULL

select`csv_to_json`(''); -- '[]'
```
---

### `csv_to_json_no_headers`

Converts a CSV file string to a json string. Will return `NULL` if given `NULL`, or if given a CSV where the values per row is not the same. This function treats all rows the same, returning each as an array of strings.

```sql
`csv_to_json_no_headers` ( `string` )
```

 - `` `string` ``
   - The CSV file string to be encoded.

## Examples

```sql
select`csv_to_json_no_headers`(
'Id,Name,Age\n'
'1,Ahmad,21\n'
'2,Ali,50');

--'[
--     [
--         "Id",
--         "Name",
--         "Age"
--     ],
--     [
--         "1",
--         "Ahmad",
--         "21"
--     ],
--     [
--         "2",
--         "Ali",
--         "50"
--     ]
-- ]'

select`csv_to_json_no_headers`(
'Id,Name,Name\n' -- Notice the double "Name" column
'1,Ahmad,21\n'
'2,Ali,50');

--'[
--     [
--         "Id",
--         "Name",
--         "Name" <-- This method doesn't change the names, unlike `csv_to_json`
--     ],
--     [
--         "1",
--         "Ahmad",
--         "21"
--     ],
--     [
--         "2",
--         "Ali",
--         "50"
--     ]
-- ]'

select`csv_to_json_no_headers`(
'Id,Name,Age,Extra\n' -- Notice this row has more values than the others
'1,Ahmad,21\n'
'2,Ali,50');                        -- NULL

select`csv_to_json_no_headers`(''); -- '[]'
```
---
## CSV Parsing Errors

As you can see in the examples and the descriptions above, both functions return `NULL` if the CSV has an error while parsing. If you do want to see the actual error messages, they will appear in the MySQL error log file like so:

```shell
tail -f /var/log/mysql/error.log # or wherever your MySQL error log is
csv-to-json: 2020/09/24 16:37:57.140331 /home/brian/go/src/github.com/StirlingMarketingGroup/mysql-csv-to-json/main.go:160: record on line 2: wrong number of fields
```

---

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install. You will also need the MySQL dev library.

Debian / Ubuntu
```shell
sudo apt update
sudo apt install libmysqlclient-dev
```

## Installing

You can find your MySQL plugin directory by running this MySQL query

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-csv-to-json.git
cd mysql-csv-to-json
go get -d ./...
go build -buildmode=c-shared -o csv_to_json.so
sudo cp csv_to_json.so /usr/lib/mysql/plugin/ # replace plugin dir here if needed
```

Enable the functions in MySQL by running this MySQL query

```sql
create function`csv_to_json`returns string soname'csv_to_json.so';
create function`csv_to_json_no_headers`returns string soname'csv_to_json.so';
```