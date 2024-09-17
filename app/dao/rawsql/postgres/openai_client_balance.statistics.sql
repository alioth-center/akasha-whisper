select ocb.client_id as client_id, oc.description as client_name, date_trunc('day', ocb.created_at) as date_day, sum(ocb.balance_change_amount) as total_cost, count(ocb.id) as request_count from openai_client_balance ocb left join openai_clients oc on ocb.client_id = oc.id where action = 'consumption' and ocb.created_at > ? group by 1, 2, 3 order by 3 desc, 1
