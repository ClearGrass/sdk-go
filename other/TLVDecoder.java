package other; 

// FIXME rename package

import java.util.*;

public class TLVDecoder {

    // 类：SubPack，表示子包
    public static class SubPack {
        public String key;
        public int len;
        public byte[] payload;

        public SubPack(String key, int len, byte[] payload) {
            this.key = key;
            this.len = len;
            this.payload = payload;
        }

        @Override
        public String toString() {
            return "{" +
                    "key='" + key + '\'' +
                    ", len=" + len +
                    ", payload=" + bytesToHex(payload) +
                    '}';
        }
    }

    // 类：TlvUnpackResult，表示 TLV 解包结果
    public static class TlvSubPackList {
        public String cmd;
        public int length;
        public List<SubPack> subPackList;

        public TlvSubPackList(String cmd, int length, List<SubPack> subPackList) {
            this.cmd = cmd;
            this.length = length;
            this.subPackList = subPackList;
        }


        @Override
        public String toString() {
            return "{" +
                    "cmd='" + cmd + '\'' +
                    ", length=" + length +
                    ", subPackList=" + subPackList +
                    '}';
        }
    }
    // 类：SensorData，表示解码后的传感器数据
    public static class SensorData {
        public String dataType;
        public int timestamp;
        public double temperature;
        public double humidity;
        public int pressure;
        public int co2;
        public int battery;
        public int rssi;
        

        @Override
        public String toString() {
            return "{" +
                    "dataType='" + dataType + '\'' +
                    ", timestamp=" + timestamp +
                    ", temperature=" + temperature +
                    ", humidity=" + humidity +
                    ", pressure=" + pressure +
                    ", co2=" + co2 +
                    ", battery=" + battery +
                    ", rssi=" + rssi +
                    '}';
        }
    }
    

    public static class TlvUnpackResult {
        public String cmd;
        public int length;
        public List<SensorData> sensorData;

        public TlvUnpackResult(String cmd, int length) {
            this.cmd = cmd;
            this.length = length;
        }


        @Override
        public String toString() {
            return "{" +
                    "cmd='" + cmd + '\'' +
                    ", length=" + length +
                    ", sensorData=" + sensorData +
                    '}';
        }
    }



    // 方法：将十六进制字符串转换为字节数组
    public static byte[] hexStringToByteArray(String s) {
        int len = s.length();
        if (len % 2 != 0) {
            throw new IllegalArgumentException("十六进制字符串长度必须为偶数");
        }
        byte[] data = new byte[len / 2];
        for(int i = 0; i < len; i += 2){
            data[i/2] = (byte)((Character.digit(s.charAt(i), 16) << 4)
                                 + Character.digit(s.charAt(i+1), 16));
        }
        return data;
    }

    // 方法：将字节数组转换为有符号整数（小端序）
    public static int bytesToIntLittleEndian(byte[] byteArray) {
        int val = 0;
        for (int i = 0; i < byteArray.length; i++) {
            val |= (byteArray[i] & 0xFF) << (i * 8);
        }
        return val;
    }

    // 方法：将字节数组转换为十六进制字符串
    public static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for(byte b : bytes){
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    // 方法：解包 TLV 数据
    public static TlvSubPackList tlvUnpack(byte[] byteArray) {
        if (byteArray.length < 5) {
            throw new IllegalArgumentException("字节数组长度不足以解包 TLV 数据");
        }

        String cmd = String.format("%02x", byteArray[2]); // byteArray[2:3].hex()
        int length = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 3, 5)); // byteArray[3:5]

        if (byteArray.length < 5 + length) {
            throw new IllegalArgumentException("字节数组长度不足以提取 payload");
        }

        byte[] payload = Arrays.copyOfRange(byteArray, 5, 5 + length);

        int index = 0;
        List<SubPack> subPackList = new ArrayList<>();
        while (index < length) {
            if (index + 1 > payload.length) {
                throw new IllegalArgumentException("子包格式错误：无法提取 key");
            }
            String key = String.format("%02x", payload[index]); // payload[index:index+1].hex()
            index += 1;

            if (index + 2 > payload.length) {
                throw new IllegalArgumentException("子包格式错误：无法提取 subLen");
            }
            int subLen = bytesToIntLittleEndian(Arrays.copyOfRange(payload, index, index + 2)); // payload[index:index+2]
            index += 2;

            if (index + subLen > payload.length) {
                throw new IllegalArgumentException("子包格式错误：subPayload 超出范围");
            }
            byte[] subPayload = Arrays.copyOfRange(payload, index, index + subLen);
            index += subLen;

            SubPack subPack = new SubPack(key, subLen, subPayload);
            subPackList.add(subPack);
        }

        return new TlvSubPackList(cmd, length, subPackList);
    }

    
    public static SensorData decodeTHData(byte[] byteArray) {
        int th = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 0, 3));
        double temperature = ((th >> 12) - 500) / 10.0;
        double humidity = (th & 0xFFF) / 10.0;
        int pressure = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 3, 5));
        int battery = byteArray[5] & 0xFF; // 无符号
 

        SensorData sensorData = new SensorData();
        sensorData.dataType = "event";
        sensorData.temperature = temperature;
        sensorData.humidity = humidity;
        sensorData.pressure = pressure;
        sensorData.battery = battery;

        return sensorData;
    }

    // 方法：解码实时数据
    public static List<SensorData> decodeRealTimeData(byte[] byteArray) {
        if (byteArray.length < 11) {
            throw new IllegalArgumentException("实时数据字节数组长度不足");
        }

        List<SensorData> sensorDataList = new ArrayList<>();

        int timestamp = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 0, 4));
        SensorData sensorData = decodeTHData(Arrays.copyOfRange(byteArray, 4, byteArray.length));
        int rssi = byteArray[10] & 0xFF; // 无符号
        if (rssi >= 128) {
            rssi -= 256;
        }

        sensorData.dataType = "event";
        sensorData.timestamp = timestamp;
        sensorData.rssi = rssi;
        sensorDataList.add(sensorData);
        return sensorDataList;
    }


    // 方法：解码历史数据
    public static List<SensorData> decodeHistoryData(byte[] byteArray) {
        List<SensorData> sensorDataList = new ArrayList<>();

        int timestamp = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 0, 4));
        int duration = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 4, 6));
        int index = 6;
        int i = 0;
        int packLen = 6;

        while (index < byteArray.length) {
            SensorData sensorData = decodeTHData(Arrays.copyOfRange(byteArray, index, index+packLen));
            sensorData.timestamp = timestamp + duration* i;
            sensorData.dataType = "data";
            sensorDataList.add(sensorData);

            index += packLen;
            i++;
        }

        return sensorDataList;
    }

     // 方法：解码v2版本的数据
     public static SensorData decodeHistoryDataV2(byte[] byteArray) {
        SensorData sensorData = new SensorData();

        int timestamp = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 0, 4));
        sensorData.timestamp = timestamp;
        
        int temperatureVal;
        int humidityVal;
        int co2Val;
        int pressureVl;

	    switch (byteArray[4]) {
            case 1:
                temperatureVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 5, 7));
                humidityVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 7, 9));
                sensorData.temperature =  temperatureVal/10.0;
                sensorData.humidity = humidityVal/10.0;
                break;
            case 2:
                temperatureVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 5, 7));
                sensorData.temperature =  temperatureVal/10.0;
                break;
            case 3:
                temperatureVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 5, 7));
                humidityVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 7, 9));
                pressureVl = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 9, 11));
                sensorData.temperature =  temperatureVal/10.0;
                sensorData.humidity = humidityVal/10.0;
                sensorData.pressure = pressureVl;
                break;
            case 4: 
                temperatureVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 5, 7));
                humidityVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 7, 9));
                co2Val = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 9, 11));
                sensorData.temperature =  temperatureVal/10.0;
                sensorData.humidity = humidityVal/10.0;
                sensorData.co2 = co2Val;
                break;
            default:
                temperatureVal = bytesToIntLittleEndian(Arrays.copyOfRange(byteArray, 5, 7));
                sensorData.temperature =  temperatureVal/10.0;
                break;
        }

       
        return sensorData;
    }

    // 方法：解析 TLV 数据
    public static TlvUnpackResult tlvDecode(byte[] byteArray) {
        TlvSubPackList subPackRet = tlvUnpack(byteArray);
        TlvUnpackResult unPackRet = new TlvUnpackResult(subPackRet.cmd,subPackRet.length);
        unPackRet.sensorData = new ArrayList<>();

        // 兼容新协议
        int batteryVal = -1;

        for (SubPack subPack : subPackRet.subPackList) {
            switch (subPack.key) {
                case "14":
                    List<SensorData> realtimeData = decodeRealTimeData(subPack.payload);
                    unPackRet.sensorData = realtimeData;
                    break;
                
                case "03":
                    List<SensorData> historyData = decodeHistoryData(subPack.payload);
                    unPackRet.sensorData = historyData;
                    break;
                
                // 下面是v2版本的解析
                case "85":
                    SensorData unitData = decodeHistoryDataV2(subPack.payload);
                    unPackRet.sensorData.add(unitData);
                    break;

                case "64":
                    batteryVal = subPack.payload[0];
                    break;
                default:
                    break;
            }
        }

        for (int i = 0; i < unPackRet.sensorData.size(); i++) {
            SensorData unitData =  unPackRet.sensorData.get(i);
            if (batteryVal >=0) {
                unitData.battery = batteryVal;
            }
        }

        return unPackRet;
    }

    // 主方法
    public static void main(String[] args) {
        String src = "43474235003802005b00650100ba6401005c110500312e322e31850b00288b11680413019901e101850b00309211680419019101df011d010001470b";
        byte[] bs = hexStringToByteArray(src);
        TlvUnpackResult unpackData = tlvDecode(bs);
        System.out.println(unpackData);

        // String src = "Q0cxMAA4AgApAB0BAAEDJAAc4j9nhAOsQS4AADarMS4AADatMS4AADauMS4AADauMS4AADYYCg==";
        // byte[] bs = Base64.getDecoder().decode(src);
        // TlvUnpackResult unpackData = tlvDecode(bs);
        // System.out.println(unpackData);
    }
}
