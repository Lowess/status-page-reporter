# :fire: Status-Page-Reporter

Simple go program to track number of Incidents from a given Bitbucket status pag  given a time range and plot them as a heatmap.

## :arrow_heading_down: Installation

`go install github.com/Lowess/status-page-reporter@latest`

Or grab the binary on the release page: [status-page-reporter](https://github.com/Lowess/status-page-reporter/releases/download/v1.0.0/status-page-reporter)


## :pencil2: Usage

```
./status-page-reporter --help

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

## Tips

```sh
# Define the FROM and TO dates based on the current quarter
export FROM=$(date -j -v-3m -f "%Y-%m-%d" "$(date +'%Y')-$(printf '%02d' $(((($(date +%m)-1)/3)*3+1)))-01" +'%Y-%m-%d')
# Calculate the end date of the previous quarter
export TO=$(date -j -v1d -v-1d -v-3m -v+3m -f "%Y-%m-%d" "$FROM" +'%Y-%m-%d')

status-page-reporter --from "${FROM}" --to "${TO}" --output png > heatmap.png
status-page-reporter --from "${FROM}" --to "${TO}" --output json | jq  'keys | length'
status-page-reporter --from "${FROM}" --to "${TO}" --output json | jq 'add'
```

---

## :pray: Credits

Big thanks to [@nikolaydubina](https://github.com/nikolaydubina) for all his great work on [calendarheatmap](https://calendarheatmap.io/) which is used in this project
