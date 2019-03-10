# Visualizing

## DataPackages pre-reqs 

Install DataPackages reader for Pandas:
```
$ install pandas-datapackage-reader
```

## Altair
Use Altair to load datapackages + do simple visualization:

* Use pipenv to setup a VirtualEnv for Visualization
```
$ pipenv shell
```
* Inside the pipenv:
```
(odd2019) $ pip install pandas jupyterlab vega_datasets altair
(odd2019) $ jupyter lab
```

## Loading data

Step to load data and just perform some basic visualization:
```python
from pandas_datapackage_reader import read_datapackage
proj = read_datapackage("./workspace/2-package-csv/projek-terbengkelai-digulung-tahun-2017")

proj._metadata
{'profile': 'tabular-data-resource',
 'encoding': 'utf-8',
 'schema': {'fields': [{'name': 'Bil', 'type': 'integer', 'format': 'default'},
   {'name': 'NEGERI', 'type': 'string', 'format': 'default'},
   {'name': 'PROJEK', 'type': 'string', 'format': 'default'},
   {'name': 'PEMAJU / PELIKUIDASI', 'type': 'string', 'format': 'default'},
   {'name': 'BIL. UNIT\n DIBINA', 'type': 'integer', 'format': 'default'},
   {'name': 'BIL. UNIT DIJUAL', 'type': 'integer', 'format': 'default'},
   {'name': 'TARIKH SEPATUT SIAP', 'type': 'string', 'format': 'default'},
   {'name': 'TARIKH TERBENGKALAI', 'type': 'string', 'format': 'default'}],
  'missingValues': ['']},
 'format': 'csv',
 'mediatype': 'text/csv',
 'dialect': {'caseSensitiveHeader': False,
  'delimiter': ',',
  'doubleQuote': True,
  'header': True,
  'lineTerminator': '\r\n',
  'quoteChar': '"',
  'skipInitialSpace': True},
 'licenses': [{'name': 'CC-BY-4.0',
   'title': 'Creative Commons Attribution 4.0',
   'path': 'https://creativecommons.org/licenses/by/4.0/'}],
 'name': 'digulung0tahun-2017',
 'path': 'data/digulung-tahun-2017.csv'}

```