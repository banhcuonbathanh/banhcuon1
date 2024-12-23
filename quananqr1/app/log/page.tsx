"use client";

import React, { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronRight, LayoutList, Table2 } from "lucide-react";
import { loggerPaths, LogPath } from "@/lib/log";

// Utility functions
const getLogRange = (ids: number[]): string => {
  const min = Math.min(...ids);
  const max = Math.max(...ids);
  return `${min}-${max}`;
};

const getComponentName = (path: string): string => {
  const cleanPath = path.replace(/\.[jt]sx?$/, "");
  const fileName = cleanPath.split("/").pop() || "";
  const parts = fileName.split("-");
  return parts
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
};

const getPathType = (path: string): string => {
  if (path.includes("/admin/")) return "Admin";
  if (path.includes("/(client)/")) return "Client";
  if (path.includes("/form/")) return "Auth";
  if (path.includes("/web-socket/")) return "Core";
  return "Other";
};

const LoggerSummary = () => {
  const [expandedPaths, setExpandedPaths] = useState<string[]>([]);
  const [viewMode, setViewMode] = useState<"cards" | "table">("table");

  const togglePath = (path: string) => {
    setExpandedPaths((prev) =>
      prev.includes(path) ? prev.filter((p) => p !== path) : [...prev, path]
    );
  };

  const toggleView = () => {
    setViewMode((prev) => (prev === "cards" ? "table" : "cards"));
  };

  const TableView = () => (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-56">Component</TableHead>
            <TableHead>Log IDs</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Path Type</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {loggerPaths.map((config: LogPath) => (
            <TableRow key={config.path}>
              <TableCell className="font-medium">
                {getComponentName(config.path)}
              </TableCell>
              <TableCell>{getLogRange(config.enabledLogIds)}</TableCell>
              <TableCell>
                <Badge variant={config.enabled ? "default" : "secondary"}>
                  {config.enabled ? "Enabled" : "Disabled"}
                </Badge>
              </TableCell>
              <TableCell>
                <Badge variant="outline">{getPathType(config.path)}</Badge>
              </TableCell>
              <TableCell className="text-right">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => togglePath(config.path)}
                >
                  {expandedPaths.includes(config.path) ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );

  const CardView = () => (
    <div className="grid gap-6">
      {loggerPaths.map((config: LogPath) => (
        <Card key={config.path} className="overflow-hidden">
          <CardHeader
            className="cursor-pointer hover:bg-slate-50"
            onClick={() => togglePath(config.path)}
          >
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <CardTitle className="text-lg font-medium">
                  {getComponentName(config.path)}
                </CardTitle>
                <p className="text-sm text-gray-500">{config.path}</p>
              </div>
              <div className="flex items-center gap-3">
                <Badge variant="outline">{getPathType(config.path)}</Badge>
                <Badge variant={config.enabled ? "default" : "secondary"}>
                  {config.enabled ? "Enabled" : "Disabled"}
                </Badge>
                {expandedPaths.includes(config.path) ? (
                  <ChevronDown className="h-5 w-5" />
                ) : (
                  <ChevronRight className="h-5 w-5" />
                )}
              </div>
            </div>
          </CardHeader>

          {expandedPaths.includes(config.path) && (
            <CardContent className="pt-4">
              <div className="space-y-4">
                <h3 className="font-medium">
                  Enabled Log IDs: {config.enabledLogIds.join(", ")}
                </h3>
                <div className="grid gap-3">
                  {Object.entries(config.logDescriptions).map(([id, log]) => (
                    <div key={id} className="p-3 rounded-lg border bg-slate-50">
                      <div className="flex items-start justify-between mb-2">
                        <span className="font-medium">Log #{id}</span>
                        <Badge
                          variant={
                            log.status === "enabled" ? "default" : "secondary"
                          }
                        >
                          {log.status}
                        </Badge>
                      </div>
                      <p className="text-sm text-gray-600 mb-1">
                        {log.description}
                      </p>
                      <p className="text-xs text-gray-500">
                        Location: {log.location}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </CardContent>
          )}
        </Card>
      ))}
    </div>
  );

  return (
    <div className="p-6 max-w-6xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Logger Configuration Summary</h1>
        <Button
          variant="outline"
          onClick={toggleView}
          className="flex items-center gap-2"
        >
          {viewMode === "cards" ? (
            <>
              <Table2 className="h-4 w-4" />
              Table View
            </>
          ) : (
            <>
              <LayoutList className="h-4 w-4" />
              Card View
            </>
          )}
        </Button>
      </div>

      {viewMode === "table" ? <TableView /> : <CardView />}
    </div>
  );
};

export default LoggerSummary;
