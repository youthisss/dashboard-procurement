"use client";

import { SimpleMasterList } from "@/components/master-data/simple-master-pages";

export default function UomList() {
  return (
    <SimpleMasterList
      title="UOM Master Data"
      addLabel="Add UOM"
      searchPlaceholder="Search UOM name..."
      nameTitle="UOM Name"
    />
  );
}
