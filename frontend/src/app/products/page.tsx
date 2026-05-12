"use client";

import { SimpleMasterList } from "@/components/master-data/simple-master-pages";

export default function ProductList() {
  return (
    <SimpleMasterList
      title="Product Master Data"
      addLabel="Add Product"
      searchPlaceholder="Search product name..."
      nameTitle="Product Name"
    />
  );
}
