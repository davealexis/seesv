# SeeSV

SeeSV is a library for providing fast access to data in very large delimited data files (CSV, TSV, pipe-delimited, etc) in a memory-efficient manner.

This library grew out of the frustration of reading and analyzing very large (multi-gigabyte, 5+ million row) CSV files (and other delimited file formats). There are few options for reading these files to quickly find problems, generate a filtered subset, or find specific data. There was the constant need for both a command line tool that can be used on a server via an SSH terminal connection and a desktop tool that can quickly churn through huge files. There are many tools out there, but most fall flat when handling very large files - either they take forever to open the file or (and) end up crashing by running out of memory.

The goal with SeeSV is to provide a re-usable library that encapsulates handling delimited files, that will enable building cross-platform, command line and desktop GUI applications that provide the user with super fast viewing of data files.

## Design

### SeeSV Features

- Blazing fast loading of files so that the user gets to a productive state with the data within seconds.
- Minimal memory consumption even for extremely large files.
- The ability to jump to any area of a file in O(1) time, then provide a bounded or unbounded stream of parsed records from that point.
- Handle files with or without headers, and files with extra header lines (like file summary metadata, etc) that should be ignored.
- Simple, intuitive API.


### Design

When a file is opened, SeeSV performs a number of discovery tasks:

1. Skip past unwanted rows at the top of the file.
2. Extract the column headers, if the file contains any. The headers are stored in a list accessible as a property of the DelimitedFile object.
3. Scans the file to generate an internal index of the byte positions of the start of every line (row) in the file, excluding the headers. This index is the only  aspect of SeeSV that may use a significant, though minimal, amount of memory. In testing, the scan of a 5 million row file (2.3GB) took around 4 - 6 seconds and produced an index using 18MB of memory. The index enables a constant-time seek to any part of a file by row number.
4. As a consequence of step 3, the row count is obtained and made available as a property of the DelimitedFile object.
5. File size is also made available through a property.


## Examples

**Example 1:** Open a normal file that has a header line:

```go
package main

import (
    "fmt"
    "davealexis/seesv"
)

func main() {
    var csvFile seesv.DelimitedFile
    err := csvFile.Open("testdata/test.csv", 0, true)
    if err != nil {
        log.Fatal("Failed to open file")
    }
    defer csvFile.File.Close()

    fmt.Println("The file has", csvFile.RowCount, "rows. Size is", csvFile.Size, "bytes.")

    fmt.Println(csvFile.Headers)

    for row := csvFile.Rows(0, -1) {
        // Do something with row
    }
}
```

**Example 2:** Similar to Example 1, except the file contains two extra metadata lines at the top before the column headers that we want to skip:

```go
package main

import (
    "fmt"
    "davealexis/seesv"
)

func main() {
    var csvFile seesv.DelimitedFile
    err := csvFile.Open("testdata/test.csv", 2, true)
    if err != nil {
        log.Fatal("Failed to open file")
    }
    defer csvFile.File.Close()

    fmt.Println(csvFile.Headers)

    // Get 1,000 rows starting from row 25,000
    for row := csvFile.Rows(25_000, 1_000) {
        // Do something with row
    }
}
```

**Example 3:** The file does not contain a header row, so we just want access to the data:

```go
package main

import (
    "fmt"
    "davealexis/seesv"
)

func main() {
    var csvFile seesv.DelimitedFile
    err := csvFile.Open("testdata/test.csv", 0, false)
    if err != nil {
        log.Fatal("Failed to open file")
    }
    defer csvFile.File.Close()

    for row := csvFile.Rows(0, -1) {
        // Do something with row
    }
}
```



**Example 4:** We just want to get a single row from the file:

```go
package main

import (
    "fmt"
    "davealexis/seesv"
)

func main() {
    var csvFile seesv.DelimitedFile
    err := csvFile.Open("testdata/test.csv", 0, false)
    if err != nil {
        log.Fatal("Failed to open file")
    }
    defer csvFile.File.Close()

    // Get row 10
    row := csvFile.Row(10) {
    // Do something with row
}
```

## Roadmap

The following are some of the features that are coming:

- Automatically detect column data types
- Allow user to supply column schema.
- Support for file formats other than CSV:
  - Tab-separated
  - Pipe-delimited
  - JSON (?)
  - Compressed files (e.g. myfile.csv.gz)
- Auto-detect which line contains headers (e.g. ignore any metadata rows at the top of the file)
- Filters
- Projections - get specified columns instead of all columns

### TODO:

- [ ] Column map - csfVile.Columns['NAME'].ID |.Type|etc
- [ ] Progress and error reporter channels
- [ ] Pipe delimited files
- [ ] Tab delimited files
- [ ] Fixed length files
