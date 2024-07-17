# :fire: Status-Page-Reporter

Simple go program to track number of Incidents from a given Bitbucket status pag  given a time range and plot them as a heatmap.

## :arrow_heading_down: Installation

`go install github.com/Lowess/status-page-reporter@latest`

Or grab the binary on the release page: [status-page-reporter](https://github.com/Lowess/drone-release-tracker/releases/download/v1.0.0/status-pager-reporter)


## :pencil2: Usage

```
./status-pager-reporter --help

Usage of ./status-page-reporter:
  -endpoint string
        Status page endpoint to scrape incidents from (default "https://status.verity.gumgum.com")
  -from string
        Releases after this date will be included (default "2024-07-01")
  -output string
        Output format (json, png, jpeg, gif, svg) (default "png")
  -to string
        Releases before this date will be included (default "2024-09-30")
```

---

## :pray: Credits

Big thanks to [@nikolaydubina](https://github.com/nikolaydubina) for all his great work on [calendarheatmap](https://calendarheatmap.io/) which is used in this project
