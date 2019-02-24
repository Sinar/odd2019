# DataPackages

- Store packaged datapackages here 

## Structure

A DataPackage should be an independent folder with the following structure; as it will be zipped before being submitted to a Data Portal of any sort e.g Sinar Project's own Open Data Portal - https://data.sinarproject.org

```
<descriptive-name-of-dataset>
     |--- data <-- Cleaned, transformed csv data, saved via DataCurator or similar tool
     |--- scripts <-- Any scripts used to get original data
     |--- docs <--
     datapackage.json <-- Created using DataCurator or similar tool
     README.md <-- Document original source of data/provenance/contacts; e.g. copy details from the data.gov.my page

```