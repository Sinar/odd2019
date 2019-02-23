# Converting XLSX

A lot of the files will become available in XLSX format.

### Single file processing

```
$ gocsv xlsx raw-datasets/ncr-sabah.xlsx
```

### Lazy way

```
$ for i in `ls raw-datasets/*.xlsx`; do gocsv xlsx "$i"; done
```

The files will be generated as a subfolder matching the name of the original XLSX folder.

e.g. The sheet name is NCR SABAH; so the CSV is in raw-datasets/ncr-sabah/NCR Sabah.csv
```
$ ls -l raw-datasets/ncr-sabah*

$ ls raw-datasets/ncr-sabah*
raw-datasets/ncr-sabah.xlsx

raw-datasets/ncr-sabah:
NCR SABAH.csv

```