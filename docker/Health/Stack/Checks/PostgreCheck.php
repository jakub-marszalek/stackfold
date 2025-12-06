<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Dotenv\Dotenv;
use Health\Stack\Contracts\HealthCheckInterface;
use PDO;
use PDOException;

class PostgreCheck implements HealthCheckInterface
{
    private const string SOURCE_DIRECTORY = '/app';
    private const string DATABASE_HOST = 'db';
    private const int DATABASE_PORT = 5432;
    private const string DEFAULT_DATABASE_NAME = 'postgres';
    private const string DEFAULT_DATABASE_USER = 'postgres';
    private const string DEFAULT_DATABASE_PASSWORD = 'Admin123';

    /**
     * @var array
     */
    private array $result = [
        'status' => 'disconnected',
        'version' => 'unknown',
        'dsn' => null,
        'all_passed' => false,
        'error' => null
    ];

    /**
     * @var PDO|null
     */
    private ?PDO $pdo = null;

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
     * Construct and load environment variables from app .env file
     */
    public function __construct()
    {
        $envDir = dirname(__DIR__, 4) . self::SOURCE_DIRECTORY;
        $envPath = $envDir . '/.env';

        if (file_exists($envPath)) {
            $dotenv = Dotenv::createImmutable($envDir);
            $dotenv->safeLoad();
        } else {
            $this->result['error'] = 'Database configuration missing (.env file not found).';
            $this->envNotFound = true;
        }

        $this->host     = self::DATABASE_HOST;
        $this->port     = self::DATABASE_PORT;
        $this->dbname   = $_ENV['DATABASE_NAME'] ?? self::DEFAULT_DATABASE_NAME;
        $this->user     = $_ENV['DATABASE_USER'] ?? self::DEFAULT_DATABASE_USER;
        $this->password = $_ENV['DATABASE_PASSWORD'] ?? self::DEFAULT_DATABASE_PASSWORD;
    }

    /**
     * Check database connection
     *
     * @return array
     */
    public function check(): array
    {
        if ($this->envNotFound) {
            return $this->result;
        }

        $this->connect();
        if ($this->pdo) {
            $this->result['status'] = 'connected';
            $this->result['version'] = $this->pdo->query('SELECT version()')->fetchColumn();
            $this->result['dsn'] =
                'pgsql://' . $this->user .
                ':' . $this->password .
                '@' . $this->host .
                ':' . $this->port .
                '/' . $this->dbname;
            $this->result['all_passed'] = true;
            $this->result['error'] = null;
        }
        return $this->result;
    }

    /**
     * Check if all checks passed
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        if (!$this->result['all_passed']) {
            return false;
        }
        return true;
    }

    /**
     * Connect to PostgreSQL database via PDO
     *
     * @return void
     */
    private function connect(): void
    {
        try {
            $this->pdo = new PDO(
                'pgsql:host=' . $this->host . ';port=' . $this->port . ';dbname=' . $this->dbname,
                $this->user,
                $this->password,
                [
                    PDO::ATTR_TIMEOUT => 5,
                    PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION
                ]
            );
        } catch (PDOException $e) {
            $this->result['error'] = $e->getMessage();
            $this->result['all_passed'] = false;
            $this->pdo = null;
        }
    }
}
