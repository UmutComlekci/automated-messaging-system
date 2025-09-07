CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    content VARCHAR(160) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'sent', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP NULL,
    external_message_id VARCHAR(255) NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_phone ON messages(phone_number);

INSERT INTO messages (id, content, phone_number, status) VALUES
('7d275de9-b8e8-4479-9c2f-af1b98d840bf', 'Hello! This is a test message.', '+1234567890', 'pending'),
('5e0c7a7f-0a54-4c94-bfc2-32aa295dd57d', 'Welcome to our service!', '+1987654321', 'pending'),
('757bc0cc-2198-44cf-ba49-dacd2349cdf4', 'Your order has been confirmed.', '+1122334455', 'pending'),
('4041e0f0-c0e5-4fdb-8fef-6ce5ca8406bb', 'Thank you for your purchase!', '+1555666777', 'pending'),
('25b2ef02-12cb-405a-8657-8d813041396b', 'Reminder: Meeting at 3 PM today.', '+1999888777', 'pending');
