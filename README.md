# sdk-go
Qingping openapi and sensor data decode

## Project structure
```
| - encoding 
    | - lora_robb.go
    | - lora_test.go
    | - tlv_decode.go
| - example
    | - encoding/lora_robb.go
| - other
    | - tlv_decode.py
    | - TLVDocode.java
| - openapi TODO
| - mqtt TODO
```

## How to use
- [Qingping Indoor Environment Monitor Lora](https://github.com/ClearGrass/sdk-go/blob/main/example/encoding/lora_robb.go)
- [Qingping Indoor Environment Monitor WiFi](https://github.com/ClearGrass/sdk-go/blob/main/example/encoding/robb.go)
- [Qingping Indoor Environment Monitor NB-IOT](https://github.com/ClearGrass/sdk-go/blob/main/example/encoding/robb.go)

## java example
```bash
java TLVDocode.java
```
## python example
```bash
python tlv_decode.py
```