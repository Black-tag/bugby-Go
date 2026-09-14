-- +goose Up
ALTER TABLE bugs 
ADD COLUMN Assigned_to UUID, 
ADD COLUMN Assigned_by UUID;