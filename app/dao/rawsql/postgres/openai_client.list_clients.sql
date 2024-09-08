WITH latest_openai_client_balance AS
         (SELECT client_id,
                 balance_remaining,
                 ROW_NUMBER() OVER (PARTITION BY client_id ORDER BY created_at DESC) AS rn
          FROM openai_client_balance)
SELECT oc.id,
       oc.description,
       oc.api_key,
       oc.endpoint,
       oc.weight,
       ocb.balance_remaining AS balance
FROM openai_clients AS oc
         JOIN latest_openai_client_balance AS ocb
              ON ocb.client_id = oc.id AND ocb.rn = 1