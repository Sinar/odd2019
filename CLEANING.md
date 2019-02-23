# Cleaning CSV

Once the raw CSV has been extracted out from XLSX

### Single file processing
```
$ gocsv clean raw-datasets/ncr-sarawak.csv >workspace/0-clean-csv/ncr-sarawak.csv
```

### Lazy way
```
# Clean top level CSVs
$ cd raw-datasets ; for i in *.csv; do  gocsv clean "$i" >"../workspace/0-clean-csv/$i" ; done ; cd -

# Clean those converted from XLSX
$ cd raw-datasets ; for j in `find *  -type d`; do for i in $j/*.csv; do mkdir -p "../workspace/0-clean-csv/$j"; gocsv clean "$i" >"../workspace/0-clean-csv/$i" ; done ; done; cd -
```

At this point; you can also run some basic stats, dimensioning + describing of data.

### Lazy way
```
# For top level CSVs
$ for i in workspace/0-clean-csv/*.csv; do ; echo "$i"; gocsv dims  "$i" ; done
workspace/0-clean-csv/12statistik-kes-jenayah-mengikut-mahkamah-jan-jun-2018.csv
Dimensions:
  Rows: 30
  Columns: 9
workspace/0-clean-csv/2017--keluasan-guna-tanah-kemudahan-awam-negeri-pp.csv
Dimensions:
  Rows: 11
...

# For those converted from XLSX
$ for i in workspace/0-clean-csv/*/*.csv; do ; echo "$i"; gocsv dims "$i" ; done
workspace/0-clean-csv/bilangan-kes-jenayah-mengikut-kategori-tahun-2013-hingga-2017/Sheet1.csv
Dimensions:
  Rows: 9
  Columns: 10
workspace/0-clean-csv/gis---pengurusan-gunatanah-semasa-bagi-gunatanah-perniagaan-dan-perkhidmatan-pada-tahun-2010/Sheet1.csv
Dimensions:
  Rows: 740
  Columns: 7
workspace/0-clean-csv/gis---pengurusan-gunatanah-semasa-bagi-hutan-pada-tahun-2010/Sheet1.csv
Dimensions:
  Rows: 138
  Columns: 7
workspace/0-clean-csv/gis---pengurusan-gunatanah-semasa-bagi-institusi-dan-kemudahan-masyarakat-pada-tahun-2010/Sheet1.csv
Dimensions:
  Rows: 241
  Columns: 7
..
```

Remove extraneous spaces until the header is at the top:
```
$ gocsv behead -n 2 workspace/0-clean-csv/55-malaysias-export-destination/Sheet1.csv >workspace/1-hea
d-csv/55-malaysias-export-destination/Sheet1.csv
```

```
$ gocsv view workspace/1-head-csv/55-malaysias-export-destination/Sheet1.csv | less
----------------------+----------------------------------------+--------------+--------------------+
| Year       | Month | Country                                | Country Code | Total Export (USD) |
+-------------+-------+----------------------------------------+--------------+--------------------+
| 2016       | 1     | AFGHANISTAN                            | AF           | 6533359            |
+--------------------------------------------------------------+--------------+--------------------+
| 2016       | 1     | ALBANIA                                | AL           | 95077              |
+--------------------------------------------------------------+--------------+--------------------+
| 2016       | 1     | ALGERIA                                | DZ           | 11758400           |
...

```