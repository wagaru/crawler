# gcloud functions
crawler for two websites:
1. "https://www.amazon.com/-/zh_TW/Google-WiFi-Mesh-%E7%B3%BB%E7%B5%B1-%E8%B7%AF%E7%94%B1%E5%99%A8%E5%8D%87%E7%B4%9A/dp/B08GG9CMLR/ref=sr_1_3?keywords=google+wifi+2020&qid=1637225228&qsid=139-3079787-0386448&sr=8-3&sres=B08GGBZNZ5%2CB08GG9CMLR%2CB07CZPQ9SV%2CB07YMJ57MB%2CB01MDJ0HVG%2CB07YMKD6SM%2CB0842VS9MM%2CB08MVDLVMT%2CB08M38C3LZ%2CB08M478JZF%2CB08P7ZTVLV%2CB07H7RQHF7%2CB08F1YPMS1%2CB081VXXC44%2CB073ZMDMKH%2CB08KTQ1J9H&srpt=NETWORKING_ROUTER"
2. "https://www.apple.com/tw/shop/buy-mac/macbook-pro/14-%E5%90%8B-%E5%A4%AA%E7%A9%BA%E7%81%B0%E8%89%B2-8-%E6%A0%B8%E5%BF%83-cpu-14-%E6%A0%B8%E5%BF%83-gpu-512gb#"
3. use sendGrid api to send notification email

# workflow

gcp scheduler -> publish to topic -> gcloud pub/sub events -> trigger cloud functions

## local development

1. setup
* https://cloud.google.com/functions/docs/running/function-frameworks
* 安裝 `github.com/GoogleCloudPlatform/functions-framework-go`


## local test

在 8080 啟動 http server
```
    cd cmd/
    go run .
```

透過 curl 測試 pub/sub trigger
```
curl localhost:8080 \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "context": {
        "eventId":"1144231683168617",
        "timestamp":"2020-05-06T07:33:34.556Z",
        "eventType":"google.pubsub.topic.publish",
        "resource":{
            "service":"pubsub.googleapis.com",
            "name":"projects/sample-project/topics/gcf-test",
            "type":"type.googleapis.com/google.pubsub.v1.PubsubMessage"
        }
        },
        "data": {
        "@type": "type.googleapis.com/google.pubsub.v1.PubsubMessage",
        "attributes": {
            "attr1":"attr1-value"
        },
        "data": "d29ybGQ="
        }
    }'
```

## local deploy

* first time
```
gcloud functions deploy crawler --entry-point ParseAll --runtime go116 --trigger-topic mbp --set-env-vars SENDGRID_API_KEY=[SENDGRID_API_KEY]
```

* update
```
gcloud functions deploy crawler --entry-point ParseAll --trigger-topic mbp --set-env-vars SENDGRID_API_KEY=[SENDGRID_API_KEY]
```

# problems

* 在 local 環境抓的到 amazon 的資料，但在 gcp 上跑結果確不一樣？

# ref
* [gcloud functions deploy syntax](https://cloud.google.com/sdk/gcloud/reference/functions/deploy#--trigger-topic)