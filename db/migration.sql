CREATE TABLE IF NOT EXISTS priorities (
    id INTEGER PRIMARY KEY,         -- 1: 高, 2: 中, 3: 低, 4:最低
    code TEXT NOT NULL UNIQUE,      -- 識別子: 'P1', 'P2', 'P3', 'P4'
    emoji TEXT                      -- 表示用絵文字: 🔴🟡🟢🔵
);
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    priority_id INTEGER NOT NULL DEFAULT 4, -- なんの指定もなければ4にする
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (priority_id) REFERENCES priorities(id)
);
INSERT INTO priorities (id, code, emoji) VALUES
    (1, 'P1', '🔴'),
    (2, 'P2', '🟡'),
    (3, 'P3', '🟢'),
    (4, 'P4', '🔵')
ON CONFLICT DO NOTHING;

--- 安全性チェックと実装を後でやる