-- +goose Up
ALTER TABLE bugs 
ADD COLUMN Assigned_t0 UUID, 
ADD COLUMN Assigned_by UUID;