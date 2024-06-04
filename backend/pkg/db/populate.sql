INSERT INTO UserAccount VALUES (
    'd44e4283-c0c2-4210-8a70-2159fc7ce9e6'::uuid,
    'abc@gmail.com',
    crypt('refo', gen_salt('bf')),
    current_timestamp
) ON CONFLICT DO NOTHING;
INSERT INTO Config VALUES
('text-text-model', 'phi3'),
('image-text-model', 'llava-phi3') ON CONFLICT DO NOTHING;
INSERT INTO Story VALUES (
    '51eecb74-bd12-40b4-bd3d-71eaa2a7d71b',
    'd44e4283-c0c2-4210-8a70-2159fc7ce9e6',
    current_timestamp,
    '{}'::int[]
) ON CONFLICT DO NOTHING;
