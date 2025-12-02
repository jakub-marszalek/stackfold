<?php

/*
 * Copyright (c) 2025. Jakub Marszałek / hesoyamcc
 * All Rights Reserved.
 */

declare(strict_types=1);

namespace Health\Stack\Checks;

use Health\Stack\Contracts\HealthCheckInterface;
use Redis;
use Throwable;

class RedisCheck implements HealthCheckInterface
{
    private const string REDIS_HOST_NAME = 'redis';
    private const int REDIS_PORT = 6379;
    private const int REDIS_TIMEOUT = 5;

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
     * Check Redis connection
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
     * Check if all checks passed
     *
     * @return bool
     */
    public function isHealthy(): bool
    {
        if (!$this->allPassed) {
            return false;
        }
        return true;
    }

    /**
     * Connect to Redis
     *
     * @return bool
     */
    private function connect(): bool
    {
        $redis = new Redis();
        try {
            $redis->connect(
                self::REDIS_HOST_NAME,
                self::REDIS_PORT,
                self::REDIS_TIMEOUT
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
     * Test Redis Ping command response
     *
     * @return bool
     */
    private function testPing(): bool
    {
        $response = $this->redis->ping();
        if ($response === true) {
            $this->result['ping_response'] = 'PONG (true)';
            return true;
        }
        return false;
    }

    /**
     * Test Redis and get server info
     *
     * @return void
     */
    private function getTestInfo(): void
    {
        if ($this->testPing()) {
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
