insert into apps (id, name, secret) values (1, 'test', 'test-secret')
on conflict do nothing;