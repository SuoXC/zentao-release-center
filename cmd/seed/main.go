package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		home, _ := os.UserHomeDir()
		dbPath = filepath.Join(home, ".zentao-release-center", "release.db")
	}

	os.MkdirAll(filepath.Dir(dbPath), 0755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrate(db); err != nil {
		log.Fatal("migrate:", err)
	}

	now := time.Now().Format(time.DateTime)

	// Projects
	projects := []struct {
		id, name, desc, prodName string
		prodID                    int
	}{
		{"p1", "电商平台", "电商后台管理系统", "电商平台 V3.0", 100},
		{"p2", "用户中心", "统一用户认证与权限管理", "用户中心服务", 200},
		{"p3", "运维监控", "基础设施监控与告警平台", "监控平台 V2.0", 300},
	}
	for _, p := range projects {
		_, err := db.Exec(`INSERT OR REPLACE INTO projects
			(id, name, description, zentao_product_id, zentao_project_id, zentao_product_name, zentao_project_name, zentao_server, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, 0, ?, '', '', 'active', ?, ?)`,
			p.id, p.name, p.desc, p.prodID, p.prodName, now, now)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("项目: %s\n", p.name)
	}

	// Releases
	releases := []struct {
		id, projectID, name, version, status, summary string
	}{
		{"r1", "p1", "v3.2.0 支付模块优化", "3.2.0", "published", "优化支付流程，修复已知Bug"},
		{"r2", "p1", "v3.3.0 购物车重构", "3.3.0", "draft", "重构购物车核心逻辑，提升性能"},
		{"r3", "p2", "v2.1.0 单点登录升级", "2.1.0", "published", "升级SSO协议，支持OAuth2.1"},
		{"r4", "p3", "v2.0.1 告警规则修复", "2.0.1", "published", "修复告警规则引擎多个问题"},
	}
	for _, r := range releases {
		_, err := db.Exec(`INSERT OR REPLACE INTO releases
			(id, project_id, name, version, status, summary, publish_count, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?)`,
			r.id, r.projectID, r.name, r.version, r.status, r.summary, now, now)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  发布单: %s\n", r.name)
	}

	// Items (bugs, tasks, notes)
	items := []struct {
		id, releaseID, itemType, title, severity, priority, status, assignedTo, noteTitle, noteContent string
		zentaoID                                                                                          int
	}{
		// r1 bugs
		{"i1", "r1", "bug", "修复微信支付回调偶发超时", "2", "1", "resolved", "张伟", "", "", 1001},
		{"i2", "r1", "bug", "修复订单金额计算精度丢失问题", "1", "1", "resolved", "李明", "", "", 1002},
		{"i3", "r1", "bug", "修复退款时库存未恢复的并发问题", "2", "2", "resolved", "王芳", "", "", 1003},
		// r1 tasks
		{"i4", "r1", "task", "对接银联云闪付支付渠道", "", "1", "done", "张伟", "", "", 2001},
		{"i5", "r1", "task", "支付结果异步通知优化", "", "2", "done", "李明", "", "", 2002},
		// r1 notes
		{"i6", "r1", "note", "", "", "", "", "", "数据库迁移注意", "需要执行: ALTER TABLE orders ADD COLUMN payment_no VARCHAR(64)", 0},
		{"i7", "r1", "note", "", "", "", "", "", "灰度策略", "先开放10%流量，观察24小时后全量发布", 0},

		// r2 bugs (draft)
		{"i8", "r2", "bug", "购物车数量修改后价格未实时更新", "3", "2", "active", "赵丽", "", "", 1004},
		{"i9", "r2", "bug", "商品失效后购物车未自动移除", "3", "3", "active", "赵丽", "", "", 1005},

		// r3 bugs
		{"i10", "r3", "bug", "修复SAML回调地址校验失败", "2", "1", "resolved", "陈刚", "", "", 1006},
		{"i11", "r3", "bug", "修复Token刷新时的竞态条件", "1", "1", "resolved", "陈刚", "", "", 1007},
		// r3 tasks
		{"i12", "r3", "task", "实现OAuth2.1授权码流程", "", "1", "done", "刘洋", "", "", 2003},
		{"i13", "r3", "task", "LDAP同步用户组优化", "", "2", "done", "刘洋", "", "", 2004},

		// r4 bugs
		{"i14", "r4", "bug", "修复CPU使用率告警阈值不生效", "1", "1", "resolved", "周杰", "", "", 1008},
		{"i15", "r4", "bug", "修复告警静默期配置解析错误", "2", "2", "resolved", "周杰", "", "", 1009},
		{"i16", "r4", "bug", "修复邮件通知模板变量未替换", "3", "2", "resolved", "吴强", "", "", 1010},
		// r4 notes
		{"i17", "r4", "note", "", "", "", "", "", "升级注意", "AlertManager 需要同步升级到 v0.27+ 版本", 0},
	}

	for idx, it := range items {
		_, err := db.Exec(`INSERT OR REPLACE INTO release_items
			(id, release_id, item_type, sort_order, zentao_id, zentao_type, title, severity, priority, status, assigned_to, resolved_by, zentao_url, steps, note_title, note_content, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, '', '', '', ?, ?, ?, ?)`,
			it.id, it.releaseID, it.itemType, idx+1,
			it.zentaoID, it.itemType, it.title, it.severity, it.priority, it.status, it.assignedTo,
			it.noteTitle, it.noteContent, now, now)
		if err != nil {
			log.Fatal(err)
		}
		label := it.itemType
		if it.title != "" {
			label = it.title
		} else if it.noteTitle != "" {
			label = it.noteTitle
		}
		fmt.Printf("    %s: %s\n", it.itemType, label)
	}

	// Deployments
	deployments := []struct {
		id, releaseID, module, address, desc string
	}{
		{"d1", "r1", "Web前端", "https://shop.example.com", "Nginx + Vue SSR"},
		{"d2", "r1", "支付服务", "https://pay.example.com:8443", "Java Spring Boot"},
		{"d3", "r1", "订单服务", "https://order.example.com:8080", "Go microservice"},
		{"d4", "r3", "认证网关", "https://auth.example.com", "Kong Gateway"},
		{"d5", "r3", "用户服务", "https://user.example.com:9090", "Java Spring Boot"},
		{"d6", "r4", "Prometheus", "https://monitor.example.com:9090", "Prometheus v2.50"},
		{"d7", "r4", "AlertManager", "https://alert.example.com:9093", "AlertManager v0.27"},
		{"d8", "r4", "Grafana", "https://grafana.example.com", "Grafana v10.3"},
	}
	for _, d := range deployments {
		_, err := db.Exec(`INSERT OR REPLACE INTO release_deployments
			(id, release_id, module_name, address, description, sort_order, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 1, ?, ?)`,
			d.id, d.releaseID, d.module, d.address, d.desc, now, now)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("    部署: %s → %s\n", d.module, d.address)
	}

	// Deployments removed

	// Snapshots (for published releases)
	_ = deployments
	snapshots := []struct {
		id, releaseID, version, content string
		bugs, tasks, notes              int
	}{
		{"s1", "r1", "3.2.0", "# v3.2.0 支付模块优化 v3.2.0\n\n## 概述\n\n优化支付流程，修复已知Bug\n\n## 部署地址\n\n| 功能模块 | 地址 | 说明 |\n|---------|------|------|\n| Web前端 | https://shop.example.com | Nginx + Vue SSR |\n| 支付服务 | https://pay.example.com:8443 | Java Spring Boot |\n| 订单服务 | https://order.example.com:8080 | Go microservice |\n\n## Bug 修复（3）\n\n| # | 标题 | 严重程度 | 优先级 | 状态 | 指派给 |\n|---|------|---------|--------|------|--------|\n| 1001 | 修复微信支付回调偶发超时 | 2 | 1 | resolved | 张伟 |\n| 1002 | 修复订单金额计算精度丢失问题 | 1 | 1 | resolved | 李明 |\n| 1003 | 修复退款时库存未恢复的并发问题 | 2 | 2 | resolved | 王芳 |\n\n---\n*Generated by zentao-release-center*", 3, 2, 2},
		{"s2", "r3", "2.1.0", "# v2.1.0 单点登录升级 v2.1.0\n\n## 概述\n\n升级SSO协议，支持OAuth2.1\n\n## 部署地址\n\n| 功能模块 | 地址 | 说明 |\n|---------|------|------|\n| 认证网关 | https://auth.example.com | Kong Gateway |\n| 用户服务 | https://user.example.com:9090 | Java Spring Boot |\n\n## Bug 修复（2）\n\n| # | 标题 | 严重程度 | 优先级 | 状态 | 指派给 |\n|---|------|---------|--------|------|--------|\n| 1006 | 修复SAML回调地址校验失败 | 2 | 1 | resolved | 陈刚 |\n| 1007 | 修复Token刷新时的竞态条件 | 1 | 1 | resolved | 陈刚 |\n\n---\n*Generated by zentao-release-center*", 2, 2, 0},
		{"s3", "r4", "2.0.1", "# v2.0.1 告警规则修复 v2.0.1\n\n## 概述\n\n修复告警规则引擎多个问题\n\n## 部署地址\n\n| 功能模块 | 地址 | 说明 |\n|---------|------|------|\n| Prometheus | https://monitor.example.com:9090 | Prometheus v2.50 |\n| AlertManager | https://alert.example.com:9093 | AlertManager v0.27 |\n| Grafana | https://grafana.example.com | Grafana v10.3 |\n\n## Bug 修复（3）\n\n| # | 标题 | 严重程度 | 优先级 | 状态 | 指派给 |\n|---|------|---------|--------|------|--------|\n| 1008 | 修复CPU使用率告警阈值不生效 | 1 | 1 | resolved | 周杰 |\n| 1009 | 修复告警静默期配置解析错误 | 2 | 2 | resolved | 周杰 |\n| 1010 | 修复邮件通知模板变量未替换 | 3 | 2 | resolved | 吴强 |\n\n---\n*Generated by zentao-release-center*", 3, 0, 1},
	}
	for _, s := range snapshots {
		_, err := db.Exec(`INSERT OR REPLACE INTO release_snapshots
			(id, release_id, version, content, item_count, bug_count, task_count, note_count, published_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			s.id, s.releaseID, s.version, s.content, s.bugs+s.tasks+s.notes, s.bugs, s.tasks, s.notes, now)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Update publish counts
	for _, rID := range []string{"r1", "r3", "r4"} {
		db.Exec("UPDATE releases SET publish_count = 1 WHERE id = ?", rID)
	}

	fmt.Println("\nDemo 数据初始化完成！")
	fmt.Printf("  项目: %d\n", len(projects))
	fmt.Printf("  发布单: %d\n", len(releases))
	fmt.Printf("  条目: %d\n", len(items))
	fmt.Printf("  快照: %d\n", len(snapshots))
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		zentao_product_id INTEGER DEFAULT 0,
		zentao_project_id INTEGER DEFAULT 0,
		zentao_product_name TEXT DEFAULT '',
		zentao_project_name TEXT DEFAULT '',
		zentao_server TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	CREATE TABLE IF NOT EXISTS releases (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		name TEXT NOT NULL,
		version TEXT DEFAULT '',
		status TEXT DEFAULT 'draft',
		summary TEXT DEFAULT '',
		publish_count INTEGER DEFAULT 0,
		first_published_at DATETIME,
		last_published_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_releases_project ON releases(project_id);
	CREATE TABLE IF NOT EXISTS release_items (
		id TEXT PRIMARY KEY,
		release_id TEXT NOT NULL,
		item_type TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		zentao_id INTEGER DEFAULT 0,
		zentao_type TEXT DEFAULT '',
		title TEXT DEFAULT '',
		severity TEXT DEFAULT '',
		priority TEXT DEFAULT '',
		status TEXT DEFAULT '',
		assigned_to TEXT DEFAULT '',
		resolved_by TEXT DEFAULT '',
		zentao_url TEXT DEFAULT '',
		steps TEXT DEFAULT '',
		note_title TEXT DEFAULT '',
		note_content TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_items_release ON release_items(release_id);
	CREATE TABLE IF NOT EXISTS release_snapshots (
		id TEXT PRIMARY KEY,
		release_id TEXT NOT NULL,
		version TEXT DEFAULT '',
		content TEXT NOT NULL,
		item_count INTEGER DEFAULT 0,
		bug_count INTEGER DEFAULT 0,
		task_count INTEGER DEFAULT 0,
		note_count INTEGER DEFAULT 0,
		published_at DATETIME NOT NULL,
		FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_snapshots_release ON release_snapshots(release_id);
	`
	_, err := db.Exec(schema)
	return err
}
