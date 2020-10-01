GOFRA
===
### (GOlang Free Restful Asterisk)

REST API for Asterisk, FreePBX

support:

* AMI
* ARI
* CDR

Get statistic

Example:
```
curl --location --request GET 'http://localhost:8053/api/cdr?StartDate=%272020-09-17%2000:00:00%27&EndDate=%272020-09-17%2023:59:59%27'
```

Download record file

Example:
```
curl --location --request GET 'http://localhost:8053/file/2020/09/11/filename.wav'
```

Inspired by https://github.com/incu6us/asterisk-ami-api
