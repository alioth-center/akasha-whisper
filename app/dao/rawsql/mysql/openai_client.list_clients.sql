SELECT oc.id,
       oc.description,
       oc.api_key,
       oc.endpoint,
       oc.weight,
       ocb.balance_remaining AS balance
FROM openai_clients AS oc
         JOIN (SELECT ob.client_id,
                      ob.balance_remaining
               FROM openai_client_balance ob
                        JOIN (SELECT client_id, MAX(created_at) AS latest_created_at
                              FROM openai_client_balance
                              GROUP BY client_id) latest
                             ON ob.client_id = latest.client_id AND ob.created_at = latest.latest_created_at) AS ocb
              ON ocb.client_id = oc.id;
