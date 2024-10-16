
# Little endian byte to integer
def bytesToIntLittleEndian(byteArray):
    val = 0;
    for i in range(len(byteArray)): 
        val = val | byteArray[i] << i*8
    return val


# 解析以4347开头的hex数据
def tlvUnpack(byteArray):
    cmd = byteArray[2:3].hex()
    length = bytesToIntLittleEndian(byteArray[3:5])
    payload = byteArray[5:5+length]

    index = 0
    subPackList = [];
    while index < length:
        key = payload[index:index+1].hex()
        subLen = bytesToIntLittleEndian(payload[index+1:index+3])
        subPayload = payload[index+3:index+3+subLen]
        index = index+3+subLen
        subPack = {
            'key': key,
            'len': subLen,
            'payload': subPayload,
        }

        #print('key:',key,'payload:',subPayload.hex())
        subPackList.append(subPack)

    return {
        'cmd': cmd,
        'length': length,
        'subPackList': subPackList,
    }

def decodeRealTimeData(byteArray):
    timestamp = bytesToIntLittleEndian(byteArray[0:4])
    th = bytesToIntLittleEndian(byteArray[4:7])
    temperature = ((th >> 12) - 500) / 10
    humidity = (th & 0xFFF) / 10
    pressure = bytesToIntLittleEndian(byteArray[7:9])
    battery = byteArray[9]
    # 有符号的整数
    rssi = byteArray[10]
    if rssi >= 128:
        rssi = rssi - 256

    return {
        'dataType': 'event',
        'timestamp': timestamp,
        'temperature': temperature,
        'humidity': humidity,
        'pressure': pressure,
        'battery': battery,
        'rssi':rssi,
    }

def decodeHistoryData(byteArray):
    pass


def tlvDecode(byteArray):
    unpackData = tlvUnpack(byteArray)

    for subPack in unpackData['subPackList']: 
        if subPack['key'] == '14':
            realtimeData = decodeRealTimeData(subPack['payload'])
            print(realtimeData)
        
        if subPack['key'] == '03':
            decodeHistoryData(subPack['payload'])
    pass

# test
if __name__ == '__main__':
    src = '43473442003802002900110500322e302e36220400303030302c01000067040004000000340500312e392e35350500322e302e361d010001140c00a82b0f6707332e00003ae6006109'
    bs = bytes.fromhex(src)
    tlvDecode(bs)


