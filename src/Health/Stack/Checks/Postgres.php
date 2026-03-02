<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Dotenv\Dotenv;
use Health\Stack\Contracts\HealthCheckInterface;
use PDO;
use PDOException;

class Postgres implements HealthCheckInterface
{
    private const string HOST = 'db';
    private const int PORT = 5432;
    private const string DEFAULT_DB = 'database';
    private const string DEFAULT_USER = 'admin';
    private const string DEFAULT_PASS = 'Admin123';

    private const string SOURCE_DIRECTORY = '/var/www/html';

    /**
     * @var array
     */
    private array $result = [
        'status' => 'disconnected',
        'version' => 'unknown',
        'masked_dsn' => null,
        'all_passed' => false,
        'error' => null
    ];

    /**
     * @var bool
     */
    private bool $envNotFound = false;
    /**
     * @var string
     */
    private string $host;
    /**
     * @var int
     */
    private int $port;
    /**
     * @var string|null
     */
    private ?string $dbname;
    /**
     * @var string|null
     */
    private ?string $user;
    /**
     * @var string|null
     */
    private ?string $password;

    /**
     * Postgres constructor.
     */
    public function __construct()
    {
        $envDir = self::SOURCE_DIRECTORY;
        $envPath = $envDir . '/.env';

        if (file_exists($envPath)) {
            $dotenv = Dotenv::createImmutable($envDir);
            $dotenv->safeLoad();
        } else {
            $this->result['error'] = 'Database configuration missing (.env file not found).';
            $this->envNotFound = true;
        }

        $this->host     = self::HOST;
        $this->port     = self::PORT;
        $this->dbname   = $_ENV['DB_DATABASE'] ?? self::DEFAULT_DB;
        $this->user     = $_ENV['DB_USERNAME'] ?? self::DEFAULT_USER;
        $this->password = $_ENV['DB_PASSWORD'] ?? self::DEFAULT_PASS;
    }

    /**
     * Check the connection with database.
     *
     * @return array
     */
    public function check(): array
    {
        if ($this->envNotFound) {
            return $this->result;
        }

        $pdo = $this->connect();
        if ($pdo instanceof PDO) {
            $this->result['status'] = 'connected';
            $this->result['version'] = $pdo->query('SELECT version()')->fetchColumn();
            $this->result['masked_dsn'] = sprintf(
                'pgsql:host=%s;port=%d;dbname=%s',
                $this->host,
                $this->port,
                $this->dbname
            );
            $this->result['all_passed'] = true;
        }
        return $this->result;
    }

    /**
     * Check if all checks were passed.
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        return $this->result['all_passed'];
    }

    /**
     * Connect to the database via PDO.
     *
     * @return PDO|null
     */
    private function connect(): ?PDO
    {
        try {
            $dsn = sprintf('pgsql:host=%s;port=%d;dbname=%s', $this->host, $this->port, $this->dbname);
            return new PDO(
                $dsn,
                $this->user,
                $this->password,
                [PDO::ATTR_TIMEOUT => 2, PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION,]
            );
        } catch (PDOException $e) {
            $this->result['error'] = 'Connection failed: ' . $e->getMessage();
            $this->result['all_passed'] = false;
            return null;
        }
    }
}
