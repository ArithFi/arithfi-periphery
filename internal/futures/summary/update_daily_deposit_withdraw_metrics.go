package summary

// 每日的充值和提现汇总数据
// 用户、总充值、总提现
/*
  SELECT walletAddress, SUM(amount) AS deposit
  FROM deposit_withdrawal
  WHERE CONVERT_TZ(_createTime, '+00:00', '+08:00') >= '2023-12-01 00:00:00'
    AND CONVERT_TZ(_createTime, '+00:00', '+08:00') < '2023-12-02 00:00:00'
    AND orderType IN ('DEPOSIT', 'WALLET_DEPOSIT')
  GROUP BY walletAddress
*/
