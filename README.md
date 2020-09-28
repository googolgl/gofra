# GOFRA 
### (GOlang Free Restful Asterisk)

REST API for Asterisk, FreePBX

support:

* AMI
* ARI
* CDR


Example:
```
curl --location --request GET 'http://localhost:8053/api/cdr?startDate=%272020-09-17%2000:00:00%27&endDate=%272020-09-17%2023:59:59%27'
```

Inspired by https://github.com/incu6us/asterisk-ami-api
