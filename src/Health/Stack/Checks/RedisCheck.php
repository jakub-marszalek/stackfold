<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;
use Redis;
use Throwable;

class RedisCheck implements HealthCheckInterface
{
    private const string HOST = 'redis';
    private const int PORT = 6379;
    private const int TIMEOUT = 2;

    /**
     * @var Redis|null
     */
    private ?Redis $redis = null;

    /**
     * @var array
     */
    private array $result = [
        'status' => 'disconnected',
        'ping_response' => null,
        'test_value' => 'failed',
        'error' => null,
        'info' => null
    ];

    /**
     * @var bool
     */
    private bool $allPassed = false;

    /**
     * Check the connection with Redis.
     *
     * @return array
     */
    public function check(): array
    {
        if ($this->connect()) {
            $this->getTestInfo();
            if ($this->result['test_value'] === 'passed') {
                $this->allPassed = true;
            }
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
        return $this->allPassed;
    }

    /**
     * Connect to the Redis.
     *
     * @return bool
     */
    private function connect(): bool
    {
        $redis = new Redis();
        try {
            $redis->connect(
                self::HOST,
                self::PORT,
                self::TIMEOUT
            );
            $this->redis = $redis;
            $this->result['status'] = 'connected';
            return true;
        } catch (Throwable $e) {
            $this->result['error'] = $e->getMessage();
            $this->redis = null;
        }
        return false;
    }

    /**
     * Check Redis response for the ping command.
     *
     * @return bool
     */
    private function checkPing(): bool
    {
        $response = $this->redis->ping();
        if ($response === true) {
            $this->result['ping_response'] = 'PONG (true)';
            return true;
        }
        return false;
    }

    /**
     * Test Redis with the ping command and get server info.
     *
     * @return void
     */
    private function getTestInfo(): void
    {
        if ($this->checkPing()) {
            $testKey = 'health_check' . time();
            $testValue = 'test_' . uniqid();

            $this->redis->set($testKey, $testValue);
            $getValue = $this->redis->get($testKey);

            $this->result['test_value'] = ($getValue === $testValue) ? 'passed' : 'failed';
            $this->result['info'] = $this->redis->info('server');

            $this->redis->del($testKey);
        }
    }
}
