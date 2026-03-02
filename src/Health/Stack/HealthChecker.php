<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack;

use Health\Stack\Contracts\HealthCheckerInterface;
use Health\Stack\Contracts\HealthCheckInterface;

class HealthChecker implements HealthCheckerInterface
{
    /**
     * @var array<string, HealthCheckInterface>
     */
    private array $checks = [];

    /**
     * @var array
     */
    private array $result = [];

    /**
     * Register a health check.
     *
     * @param string $name
     * @param HealthCheckInterface $check
     * @return void
     */
    public function register(string $name, HealthCheckInterface $check): void
    {
        $this->checks[$name] = $check;
    }

    /**
     * Run all registered health checks.
     *
     * @return array|array[]
     */
    public function run(): array
    {
        foreach ($this->checks as $name => $check) {
            $this->result[$name] = $check->check();
        }
        $this->summarize();
        return $this->result;
    }

    /**
     * Summarize results.
     *
     * @return void
     */
    public function summarize(): void
    {
        $this->result['summary']['overall_status'] = 'PASSED';

        foreach ($this->checks as $name => $check) {
            $isHealthy = $check->isHealthy();
            if ($isHealthy === true) {
                $this->result['summary'][$name] = 'OK';
            } else {
                $this->result['summary'][$name] = 'ERROR';
                $this->result['summary']['overall_status'] = 'FAILED';
            }
        }
    }
}
