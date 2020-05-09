# 比特币技术实现
Go 实现简易版比特币。

## 运行脚本

```shell
bin/
  linux   - Linux系统
  mac     - MacOS系统
  windows - Windows系统
```

## 命令

```shell
bitcoin-go\bin\windows>.\bitcoin
  -clear
        删除所有区块
  -create-block-chain string
        创建区块链
  -create-wallet
        创建钱包
  -get-balance string
        获取余额
  -list
        显示所有区块
  -list-transaction
        显示所有交易
  -list-wallet
        显示所有钱包地址
  -send
        转账（付款人 收款人 转账金额 miner）
```

## 创建钱包
每个钱包包含用于签名的私钥和公钥。

命令:

```shell
.\bitcoin -create-wallet
```
创建三个钱包:

```shell
bitcoin-go\bin\windows>.\bitcoin -create-wallet
钱包地址: 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf
bitcoin-go\bin\windows>.\bitcoin -create-wallet
钱包地址: 1Q919Bek615WSetANgGccoUgTwpp76xp8b
bitcoin-go\bin\windows>.\bitcoin -create-wallet
钱包地址: 14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc
```

## 显示所有钱包

命令:

```shell
.\bitcoin -list-wallet
```
钱包列表:

```shell
bitcoin-go\bin\windows>.\bitcoin -list-wallet
钱包地址: 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf
钱包地址: 1Q919Bek615WSetANgGccoUgTwpp76xp8b
钱包地址: 14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc
```

## 创建区块链

命令:

```shell
.\bitcoin -create-block-chain 钱包地址
```

创建区块链:

```shell
bitcoin-go\bin\windows>.\bitcoin -create-block-chain 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf
挖矿成功! nonce: 69950, hash: 00006dd322d34fd42c25ac2975a4a48f6089141c97ca32c1745f2e367607d186
创建区块链成功!!!
```

区块链创建成功后，产生一笔挖矿交易，默认奖励为 12.5。

`transaction.go`

```go
const reward = 12.5
```

## 获取余额

命令:

```shell
.\bitcoin -get-balance 钱包地址
```

查看余额:

```shell
bitcoin-go\bin\windows>.\bitcoin -get-balance 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf
1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf的余额为12.500000
```

## 转账

命令:

```shell
.\bitcoin -send 付款人 收款人 转账金额 矿工
```

`1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf` 向 `1Q919Bek615WSetANgGccoUgTwpp76xp8b` 转 2.5，指定 `14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc` 为矿工。

```shell
bitcoin-go\bin\windows>.\bitcoin -send 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf 1Q919Bek615WSetANgGccoUgTwpp76xp8b 2.5 14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc
对交易进行签名...
找到交易 311891e166c644466bf03486d5ce6ebf361f927f3d019b7fa5df4482f472261f
签名...
对数据 fbc367d09d9e6a3193f8877a5a053523b4d1b89b4464f3fda5a82da6641e2370 进行签名
对交易进行校验...
对交易进行校验...
找到交易 311891e166c644466bf03486d5ce6ebf361f927f3d019b7fa5df4482f472261f
校验...
对数据 fbc367d09d9e6a3193f8877a5a053523b4d1b89b4464f3fda5a82da6641e2370 进行校验
挖矿成功! nonce: 35610, hash: 0000e69fc5491eaaa6274d2b31942e114f70c8b7aa43191a57dcad4bb603c2f9
```

获取 `1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf` 余额:

```shell
bitcoin-go\bin\windows>.\bitcoin -get-balance 1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf
1HHqc8uJyUg7qtwbbqP9MNEpeTDSBykmrf的余额为10.000000
```

获取 `1Q919Bek615WSetANgGccoUgTwpp76xp8b` 余额:

```shell
bitcoin-go\bin\windows>.\bitcoin -get-balance 1Q919Bek615WSetANgGccoUgTwpp76xp8b
1Q919Bek615WSetANgGccoUgTwpp76xp8b的余额为2.500000
```

获取 `14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc` 余额:

```shell
bitcoin-go\bin\windows>.\bitcoin -get-balance 14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc
14sobo6du8vRW9ZPXjMSgzfGWAcaGtwFEc的余额为12.500000
```

## 显示所有区块

命令:

```shell
.\bitcoin -list
```

列出所有区块:

```shell
bitcoin-go\bin\windows>.\bitcoin -list
=================1===================
Version: 0
PrevHash: 00006dd322d34fd42c25ac2975a4a48f6089141c97ca32c1745f2e367607d186
MerKleRoot: 60a93ab78585ab188dc42ba213c1df4ea68a1056e438a860b11b031dcf4d54bd
Timestamp: 1589033835
Difficulty: 16
Nonce: 35610
Hash: 0000e69fc5491eaaa6274d2b31942e114f70c8b7aa43191a57dcad4bb603c2f9
IsValid: true
=================2===================
Version: 0
PrevHash: 00
MerKleRoot: ffac8a9b299b21dda28d7a32eb46b853826b0c802bc86cd8501cb127b3dea46a
Timestamp: 1589032964
Difficulty: 16
Nonce: 69950
Hash: 00006dd322d34fd42c25ac2975a4a48f6089141c97ca32c1745f2e367607d186
IsValid: true
```

## 显示所有交易

命令:

```shell
.\bitcoin -list-transaction
```

列出所有交易:

```shell
bitcoin-go\bin\windows>.\bitcoin -list-transaction

  Transaction 7da98ecc0835699656a851abb085ea0525f7387a5a5ee343a8d2e81d90d47907:
    Input 0:
      TxId:
      OutIndex: -1
      Signature:
      PublicKey: 476f20e58cbae59d97e993be
    Output 0:
      Value: 12.500000
      PublicKeyHash: 2a841338e1617c8fa6957768823620360bed2ff6
  Transaction f37fe9fb072bc2b090ec406611d063a13bf7384f155bf947347665b02533fbce:
    Input 0:
      TxId: 311891e166c644466bf03486d5ce6ebf361f927f3d019b7fa5df4482f472261f
      OutIndex: 0
      Signature: 590c523d016f02c47c9625f019ccaa975b026f3ff2695d7d39722254944fe2915acf5c69c28143958b978496a56d5926ebbc8c98f5db1bd4d735c0cce4baf56c
      PublicKey: 2fff81030101095075626c69634b657901ff820001030105437572766501100001015801ff840001015901ff840000000aff83050102ff8600000045ff82011963727970746f2f656c6c69707469632e703235364375727665ff870301010970323536437572766501ff88000101010b4375727665506172616d7301ff8a00000053ff890301010b4375727665506172616d7301ff8a00010701015001ff840001014e01ff840001014201ff84000102477801ff84000102477901ff8400010742697453697a6501040001044e616d65010c000000fe0108ff88ffbd01012102ffffffff00000001000000000000000000000000ffffffffffffffffffffffff012102ffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325510121025ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b0121026b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c2960121024fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f501fe02000105502d32353600000121024baf545603008cc368be798f869ee4ef0b44c56f48a3ea379f602a7e2171d6cf012102bb85ea3cbdaee2d8e2fc499b1455c76ab672c93c3f31316fa3186abfc94d9f3c00
    Output 0:
      Value: 2.500000
      PublicKeyHash: fdce56a824790378dd505c3c48f52a92f909293f
    Output 1:
      Value: 10.000000
      PublicKeyHash: b2b13df40f45f628eb7a0230a3ecf4d0513f55ed
  Transaction 311891e166c644466bf03486d5ce6ebf361f927f3d019b7fa5df4482f472261f:
    Input 0:
      TxId:
      OutIndex: -1
      Signature:
      PublicKey: 476f20e58cbae59d97e993be
    Output 0:
      Value: 12.500000
      PublicKeyHash: b2b13df40f45f628eb7a0230a3ecf4d0513f55ed
```

