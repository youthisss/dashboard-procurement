"use client";

import { SimpleMasterList } from "@/components/master-data/simple-master-pages";

export default function MotList() {
  return (
    <SimpleMasterList
      title="MOT Master Data"
      addLabel="Add MOT"
      searchPlaceholder="Search MOT name..."
      nameTitle="MOT Name"
    />
  );
}
