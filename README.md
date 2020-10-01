GOFRA
===
### (GOlang Free Restful Asterisk)

REST API for Asterisk, FreePBX

support:
* AMI
* ~ARI~
* CDR

AMI request Example:
```s
curl --location --request POST 'http://localhost:8053/api/ami/async' --header 'Content-Type: application/json' --data-raw '{"AcTion":"Originate", "Channel":"PJSIP/100", "Context":"from-internal","Exten":"0XXXXXXXXX","Priority":"1","Callerid":"100"}'
```

Get statistic Example:
```s
curl --location --request GET 'http://localhost:8053/api/cdr?StartDate=%272020-09-17%2000:00:00%27&EndDate=%272020-09-17%2023:59:59%27'
```

Download record file Example:
```s
curl --location --request GET 'http://localhost:8053/file/2020/09/11/filename.wav'
```

Inspired by https://github.com/incu6us/asterisk-ami-api
