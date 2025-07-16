CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE rooms (
                       id UUID PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       private BOOLEAN NOT NULL DEFAULT FALSE,
                       category VARCHAR(100),
                       user_count INTEGER NOT NULL DEFAULT 0,
                       description TEXT,
                       owner_id UUID NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE roles (
                       id UUID PRIMARY KEY,
                       room_id UUID NOT NULL,
                       name VARCHAR(100) NOT NULL,
                       color VARCHAR(20),
                       priority INTEGER,
                       permissions BIGINT,
                       created_at TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP NOT NULL,
                       FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE TABLE room_members (
                              room_id UUID NOT NULL,
                              user_id UUID NOT NULL,
                              role_id UUID,
                              joined_at TIMESTAMP NOT NULL,
                              PRIMARY KEY (room_id, user_id),
                              FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                              FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL
);

CREATE TABLE room_invites (
                              id UUID PRIMARY KEY,
                              room_id UUID NOT NULL,
                              invited_id UUID NOT NULL,
                              sent_by_id UUID NOT NULL,
                              status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'accepted', 'declined')),
                              sent_at TIMESTAMP NOT NULL,
                              FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                              FOREIGN KEY (invited_id) REFERENCES users(id) ON DELETE CASCADE,
                              FOREIGN KEY (sent_by_id) REFERENCES users(id) ON DELETE CASCADE
);
