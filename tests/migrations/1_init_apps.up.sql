INSERT INTO storage_file (id,
                          file_name,
                          insert_date,
                          update_date,
                          file_path,
                          file_hash)
VALUES (1,
        'test',
        '01.01.2025',
        '01.01.2025',
        '/dfsg/srfg/sgdf',
        '0x112233445566778899AABBCCDDEEFF')
ON CONFLICT DO NOTHING;

INSERT INTO storage_file (id,
                          file_name,
                          insert_date,
                          update_date,
                          file_path,
                          file_hash)
VALUES (2,
        'test2',
        '01.01.2025',
        '01.01.2025',
        '/dfsg/srfg/sgdf2',
        '0x112233445566778899AABBCCDDEEFF')
    ON CONFLICT DO NOTHING;