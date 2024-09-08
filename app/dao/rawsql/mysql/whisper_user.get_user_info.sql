SELECT wu.id,
       wu.email,
       wu.api_key,
       wu.role,
       wu.language,
       wu.allow_ips,
       wu.updated_at,
       wb.balance_remaining AS balance
FROM whisper_users AS wu
         JOIN (SELECT wub.user_id,
                      balance_remaining
               FROM whisper_user_balance wub
                        JOIN (SELECT user_id, MAX(created_at) AS latest_created_at
                              FROM whisper_user_balance
                              GROUP BY whisper_user_balance.user_id) latest
                             ON wub.user_id = latest.user_id AND wub.created_at = latest.latest_created_at) AS wb
              ON wu.id = wb.user_id
WHERE wu.id = ?;
