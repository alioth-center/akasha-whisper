WITH balance AS (SELECT user_id,
                        balance_remaining,
                        ROW_NUMBER() over (PARTITION BY user_id ORDER BY created_at DESC) AS rn
                 from whisper_user_balance)
SELECT wu.id,
       wu.email,
       wu.api_key,
       wu.role,
       wu.language,
       wu.allow_ips,
       wu.updated_at,
       wb.balance_remaining AS balance
FROM whisper_users AS wu
         JOIN balance AS wb ON wu.id = wb.user_id AND wb.rn = 1 AND wu.id = ?;