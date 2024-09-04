WITH latest_openai_client_balance AS (SELECT client_id,
                                             balance_remaining,
                                             ROW_NUMBER() OVER (PARTITION BY client_id ORDER BY created_at DESC) AS rn
                                      FROM openai_client_balance),
     latest_whisper_user_balance AS (SELECT user_id,
                                            balance_remaining,
                                            ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) AS rn
                                     FROM whisper_user_balance)
SELECT oc.id                 AS client_id,
       oc.weight             AS client_weight,
       oc.description        AS client_description,
       ocb.balance_remaining AS client_balance,
       wu.id                 AS user_id,
       wu.role               AS user_role,
       wub.balance_remaining AS user_balance,
       om.model              AS model_name,
       om.id                 AS model_id,
       om.max_tokens         AS model_max_tokens,
       om.prompt_price       AS model_prompt_price,
       om.completion_price   AS model_completion_price
FROM whisper_users AS wu
         JOIN whisper_user_permissions AS wup ON wu.id = wup.user_id AND wu.api_key = '${user_api_key}'
         JOIN openai_models AS om ON wup.model_id = om.id AND om.model = '${model_name}'
         JOIN openai_clients AS oc ON oc.id = om.client_id
         JOIN latest_openai_client_balance AS ocb ON oc.id = ocb.client_id AND ocb.rn = 1 AND ocb.balance_remaining > 0
         JOIN latest_whisper_user_balance AS wub ON wu.id = wub.user_id AND wub.rn = 1 AND wub.balance_remaining > 0;