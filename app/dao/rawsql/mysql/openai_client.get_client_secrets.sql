SELECT oc.id,
       oc.api_key,
       oc.endpoint,
       oc.weight,
       ocb.balance_remaining AS balance
FROM openai_clients AS oc
         JOIN (
    SELECT client_id,
           balance_remaining
    FROM openai_client_balance
    WHERE (client_id, created_at) IN (
        SELECT client_id, MAX(created_at)
        FROM openai_client_balance
        GROUP BY client_id
    )
) AS ocb ON ocb.client_id = oc.id
WHERE oc.id IN (${openai_client_id});
