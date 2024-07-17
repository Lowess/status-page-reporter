total-downtime:
	./status-page-reporter --from 2024-01-01 --to 2024-03-31 --output json | jq add

incident-count:
	  ./status-page-reporter --from 2024-01-01 --to 2024-03-31 | jq 'keys | length' 
