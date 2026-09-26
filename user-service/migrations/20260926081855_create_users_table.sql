-- +goose Up                                                                                                                      
-- +goose StatementBegin                                                                                                          
CREATE TABLE IF NOT EXISTS users (                                                                                                
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),                                                                                
    phone VARCHAR(20) NOT NULL UNIQUE,                                                                                            
    password_hash VARCHAR(255) NOT NULL,                                                                                          
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),                                                                                
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()                                                                                 
);                                                                                                                                
                                                                                                                                                                                                                   
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);                                                                       
-- +goose StatementEnd                                                                                                            
                                                                                                                                      
-- +goose Down                                                                                                                    
-- +goose StatementBegin                                                                                                          
DROP TABLE IF EXISTS users;                                                                                                       
-- +goose StatementEnd                