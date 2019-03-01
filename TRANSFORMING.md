# Transforming 

## Objective

More in-depth investigation on the strcuture and data validity of the dataset; possibly using SQL statements; e.g. remove/find/mark rows where Year is 0
Advanced manipulation of data; via SQL mode, joining or transformation.

## Examples

- Example #1: Join two CSVs together e.g. Combine Export/Import data per country
```
$ gocsv sql -q 'SELECT a.Year, a.Month, a.Code, a.Import, b.Export, (a.Import - b.Export) AS DIFF 
  FROM import a LEFT JOIN export b  
  WHERE a.Year = b.Year AND a.Month = b.Month AND a.Code = b.Code' 
  workspace/2-package-csv/56-malaysias-import-sources/import.csv 
  workspace/2-package-csv/55-malaysias-export-destination/export.csv 
  >workspace/3-analysis-csv/malaysia-import-export-diff.csv
```
- Example #2: Examine the distinct values of a column and see if it makes sense
```
gocsv sql -q 'SELECT DISTINCT(TAHUN_BINA) FROM MAKLUMAT_BANGUNAN_BLOK' 
  workspace/2-package-csv/9.maklumat-bangunan-dan-blok/MAKLUMAT_BANGUNAN_BLOK.csv  
  | sort -r

TAHUN_BINA
81
206
2018
2017
2016
...
2000
200
..
```

- Example #3: Look for invalid data (Schools built before 1800!)
```
$ gocsv sql -q 'SELECT NEGERI,PPD,NAMASEKOLAH,TAHUN_BINA  
  FROM MAKLUMAT_BANGUNAN_BLOK WHERE TAHUN_BINA <1800' 
  workspace/2-package-csv/9.maklumat-bangunan-dan-blok/MAKLUMAT_BANGUNAN_BLOK.csv

NEGERI,PPD,NAMASEKOLAH,TAHUN_BINA
KEDAH,PPD KUALA MUDA/YAN,SMA SUNGAI PETANI,0
KEDAH,PPD KUALA MUDA/YAN,SEKOLAH KEBANGSAAN DATARAN MUDA,200
KEDAH,PPD LANGKAWI,SEKOLAH KEBANGSAAN EWA,0
NEGERI SEMBILAN,PPD TAMPIN,SEKOLAH KEBANGSAAN SUNGAI DUA,0
PAHANG,PPD PEKAN,SEKOLAH SULTAN HAJI AHMAD SHAH,1095
PERAK,PPD BATANG PADANG,SEKOLAH KEBANGSAAN BALUN,81
PERAK,PPD KINTA UTARA,SEKOLAH KEBANGSAAN MARIAN CONVENT,157
PERLIS,JPN PERLIS,SEKOLAH KEBANGSAAN ORAN,0
PERLIS,JPN PERLIS,SEKOLAH KEBANGSAAN ORAN,0
SARAWAK,PPD BARAM,SEKOLAH KEBANGSAAN LONG LELLANG,206
SARAWAK,PPD SRI AMAN,SEKOLAH KEBANGSAAN SELANJAN,0
SARAWAK,PPD SRI AMAN,SEKOLAH KEBANGSAAN SELANJAN,0
TERENGGANU,PPD SETIU,SEKOLAH MENENGAH KEBANGSAAN TENGKU IBRAHIM,0
TERENGGANU,PPD SETIU,SEKOLAH MENENGAH KEBANGSAAN TENGKU IBRAHIM,0

```

- Example #4: Extract out a matching column of data from a larger set to a smaller subset for ease of use.
Filter only schools in the state of Perlis.
```
$ wc -l workspace/2-package-csv/9.maklumat-bangunan-dan-blok/MAKLUMAT_BANGUNAN_BLOK.csv
   74451 

$ gocsv filter -c NEGERI -i -eq "PERLIS" 
  workspace/2-package-csv/9.maklumat-bangunan-dan-blok/MAKLUMAT_BANGUNAN_BLOK.csv 
  >workspace/3-analysis-csv/perlis-maklumat-bangunan-dan-blok.csv

$ wc -l workspace/3-analysis-csv/perlis-maklumat-bangunan-dan-blok.csv
  1061 
```