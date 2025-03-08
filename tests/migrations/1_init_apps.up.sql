INSERT INTO storage_file (id,
                          file_name,
                          insert_date,
                          update_date,
                          delete_date,
                          file_path,
                          file_hash)
VALUES (1,
        'test.png',
        '2025-03-08 17:56:01',
        '2025-03-08 17:56:01',
        null,
        'storage/files/2025/03/08/8022896507_test.png',
        'f2d24fea65e85073d273a0c0631ce7379901183b401e893bc83b3fa3981ce3e9')
ON CONFLICT DO NOTHING;

INSERT INTO storage_file (id,
                          file_name,
                          insert_date,
                          update_date,
                          delete_date,
                          file_path,
                          file_hash)
VALUES (2,
        'test2.png',
        '2025-03-08 17:56:01',
        '2025-03-08 17:56:01',
        null,
        'storage/files/2025/03/08/8022896502_test2.png',
        'f2d24fea65e85073d273a0c0631ce7379901183b401e893bc83b3fa3981ce3e9')
    ON CONFLICT DO NOTHING;