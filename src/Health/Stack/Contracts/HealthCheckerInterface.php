<?php

/*
 * Copyright (c) 2026. Jakub Marszałek <marszalek.jakub@proton.me>
 * Licensed under the MIT License. See LICENSE file in the project root.
 */

declare(strict_types=1);

namespace Health\Stack\Contracts;

interface HealthCheckerInterface
{
    /**
     * Register a health check.
     *
     * @param string $name
     * @param HealthCheckInterface $check
     * @return void
     */
    public function register(string $name, HealthCheckInterface $check): void;

    /**
     * Run all registered health checks.
     *
     * @return array
     */
    public function run(): array;

    /**
     * Summarize results.
     *
     * @return void
     */
    public function summarize(): void;
}
