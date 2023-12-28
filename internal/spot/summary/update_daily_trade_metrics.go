package summary

// 链上每日的交易数据，Pancake，Buy、Sell

/*
  SELECT to_address, SUM(value) buyValue
  FROM erc20_transfer_bsc
  WHERE from_address = "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38"
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') >= '2023-12-02 00:00:00'
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') < '2023-12-03 00:00:00'
  GROUP BY to_address;
*/

/*
  SELECT from_address, SUM(value) sellValue
  FROM erc20_transfer_bsc
  WHERE to_address = "0xac4c8fabbd1b7e6a01afd87a17570bbfa28c7a38"
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') >= '2023-12-02 00:00:00'
    AND CONVERT_TZ(from_unixtime(timestamp), '+00:00', '+08:00') < '2023-12-03 00:00:00'
  GROUP BY from_address;
*/
