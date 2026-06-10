-- MySQL dump 10.13  Distrib 8.0.40, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: blog_system
-- ------------------------------------------------------
-- Server version	8.0.40

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `blog_system`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `blog_system` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `blog_system`;

--
-- Table structure for table `ai_review_records`
--

DROP TABLE IF EXISTS `ai_review_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ai_review_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `target_type` varchar(20) NOT NULL,
  `target_id` bigint unsigned NOT NULL,
  `operator_id` bigint unsigned NOT NULL,
  `risk_level` varchar(10) NOT NULL DEFAULT 'low',
  `risk_labels` longtext,
  `summary` varchar(500) DEFAULT NULL,
  `suspicious_segments` longtext,
  `suggestion` varchar(500) DEFAULT NULL,
  `model_name` varchar(60) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ai_review_target` (`target_type`,`target_id`),
  KEY `idx_ai_review_records_operator_id` (`operator_id`),
  CONSTRAINT `fk_ai_review_records_operator` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ai_review_records`
--

LOCK TABLES `ai_review_records` WRITE;
/*!40000 ALTER TABLE `ai_review_records` DISABLE KEYS */;
/*!40000 ALTER TABLE `ai_review_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_favorites`
--

DROP TABLE IF EXISTS `article_favorites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_favorites` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `article_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_article_favorite` (`article_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_favorites`
--

LOCK TABLES `article_favorites` WRITE;
/*!40000 ALTER TABLE `article_favorites` DISABLE KEYS */;
INSERT INTO `article_favorites` VALUES (2,9,6,'2026-06-01 14:08:49.478');
/*!40000 ALTER TABLE `article_favorites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_likes`
--

DROP TABLE IF EXISTS `article_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_likes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `article_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_article_like` (`article_id`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_likes`
--

LOCK TABLES `article_likes` WRITE;
/*!40000 ALTER TABLE `article_likes` DISABLE KEYS */;
INSERT INTO `article_likes` VALUES (2,7,1,'2026-06-01 12:38:50.464');
/*!40000 ALTER TABLE `article_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_reviews`
--

DROP TABLE IF EXISTS `article_reviews`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_reviews` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `article_id` bigint unsigned NOT NULL,
  `reviewer_id` bigint unsigned NOT NULL,
  `action` varchar(16) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_article_reviews_article_id` (`article_id`),
  KEY `idx_article_reviews_reviewer_id` (`reviewer_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_reviews`
--

LOCK TABLES `article_reviews` WRITE;
/*!40000 ALTER TABLE `article_reviews` DISABLE KEYS */;
INSERT INTO `article_reviews` VALUES (3,4,1,'approve','','2026-04-15 19:18:47.502'),(4,13,1,'approve','','2026-06-01 14:08:37.101');
/*!40000 ALTER TABLE `article_reviews` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `article_tags`
--

DROP TABLE IF EXISTS `article_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `article_tags` (
  `article_id` bigint unsigned NOT NULL,
  `tag_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`article_id`,`tag_id`),
  KEY `fk_article_tags_tag` (`tag_id`),
  CONSTRAINT `fk_article_tags_article` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`),
  CONSTRAINT `fk_article_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `article_tags`
--

LOCK TABLES `article_tags` WRITE;
/*!40000 ALTER TABLE `article_tags` DISABLE KEYS */;
INSERT INTO `article_tags` VALUES (4,2);
/*!40000 ALTER TABLE `article_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `articles`
--

DROP TABLE IF EXISTS `articles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `articles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(120) NOT NULL,
  `summary` varchar(300) DEFAULT NULL,
  `content` longtext NOT NULL,
  `cover_image` varchar(255) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'draft',
  `reject_reason` varchar(255) DEFAULT NULL,
  `author_id` bigint unsigned NOT NULL,
  `published_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `category_id` bigint unsigned DEFAULT NULL,
  `view_count` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_articles_status` (`status`),
  KEY `idx_articles_author_id` (`author_id`),
  KEY `idx_articles_category_id` (`category_id`),
  CONSTRAINT `fk_articles_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`),
  CONSTRAINT `fk_articles_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `articles`
--

LOCK TABLES `articles` WRITE;
/*!40000 ALTER TABLE `articles` DISABLE KEYS */;
INSERT INTO `articles` VALUES (4,'以 Codex 为中心的 AI 全栈开发工作流：从 PRD 到代码联动','本文记录了一套个人实践验证过的 AI 辅助开发工作流。核心理念是：用一个最强的模型（GPT-5.4 / Codex）统一负责规划和核心编码，用垂直优化的工具（Trae、OpenCode）处理前端和终端杂务，并建立“文档与代码双向同步”的机制，避免开发中后期的混乱。','第一阶段：规划 —— 用 Codex 生成 PRD 与技术设计\n本阶段所有产出均由 Codex（GPT-5.4） 完成。建议将生成的文档保存在项目的 /docs 目录下，纳入版本控制。\n\n步骤 1.1：生成产品需求文档（PRD）\n操作场景： 项目启动或新增一个完整功能模块（Feature）时。\n\n在 Codex 中输入以下指令（可复制）：\n\nmarkdown\n你是一位资深产品经理。请为“[在此填入功能名称，如：用户积分系统]”撰写一份结构完整的产品需求文档（PRD）。\n\n要求：\n1. 包含以下章节：背景与目标、用户故事、功能范围、业务流程图（用 Mermaid 语法描述）、非功能需求、验收标准。\n2. 语言清晰、无歧义，适合开发、测试、设计等跨职能角色阅读。\n3. 输出格式为 Markdown。\n产出物： 一份 prd-[功能名].md 文件，保存至 /docs/prd/。\n\n示例输出片段：\n\nmarkdown\n## 用户故事\n- 作为一名普通用户，我希望在完成订单后获得积分，以便在未来抵扣现金。\n- 作为一名运营人员，我希望能够配置积分获取规则，以便灵活调整营销策略。\n...\n步骤 1.2：生成技术设计文档\n操作场景： PRD 定稿后，正式开始编码前。\n\n将步骤 1.1 生成的 PRD 内容作为上下文，继续在 Codex 中指令：\n\nmarkdown\n你是一位资深全栈架构师。请基于上述 PRD，生成一份详尽的技术设计文档。\n\n要求包含：\n1. 系统架构图（使用 Mermaid 绘制）。\n2. 数据库设计：表结构 SQL 语句、字段说明、索引策略。\n3. API 接口定义：遵循 OpenAPI 3.0 规范，包含请求/响应示例。\n4. 核心模块划分与关键逻辑说明。\n5. 技术栈选型建议（如适用）。\n产出物： 一份 design-[功能名].md 文件，保存至 /docs/design/。\n\n⚠️ 重要提醒： 在让 Codex 生成上述两份文档时，建议在同一个会话（Conversation）中完成，或者手动将 PRD 文件内容附加给新会话。这能确保 GPT-5.4 充分利用上下文，保证两份文档的连贯性。\n\n💻 第二阶段：开发 —— Codex 与 Trae 的协作\n步骤 2.1：后端 / 全栈编码（Codex 执行）\n操作场景： 实现业务逻辑、数据库操作、API 接口等。\n\n最佳实践：开启“规划模式”再动手\n\n在 Codex 交互界面使用 /plan 指令（或对应 UI 按钮），让 AI 先给计划，等你审批后再写代码。这能避免 AI 直接修改文件导致不可逆的混乱。\n\n指令模板：\n\nmarkdown\n请根据 `/docs/design/design-[功能名].md` 中的设计，执行以下任务：\n\n1. 创建数据库迁移文件（基于设计文档中的 SQL）。\n2. 实现 API 接口 `[接口名称]` 的业务逻辑。\n3. 为该接口编写单元测试（遵循 AAA 模式）。\n\n在修改任何文件前，先输出一个具体的执行计划。待我回复“Approve”后再开始编码。\n关键习惯：要求文档同步更新\n\n每次让 Codex 写代码时，在指令末尾加上这句话，能彻底解决“代码改了、文档没改”的老大难问题：\n\nmarkdown\n如果在实现过程中对设计有任何偏离（例如调整了字段名、修改了接口返回结构），请**同步更新** `/docs/design/design-[功能名].md` 中的对应部分。\n步骤 2.2：前端界面实现（Trae 执行）\n操作场景： 需要快速生成前端组件、页面或样式时。\n\n工作方式：\n\n打开 Trae，确认处于 SOLO 模式。\n\n将 Codex 生成的 API 文档或接口定义（如 JSON 响应示例）复制给 Trae。\n\n指令模板：\n\nmarkdown\n请根据以下 API 返回的 JSON 结构，生成一个用户积分排行榜组件。\n\n要求：\n- 使用 Vue 3 + TypeScript + Tailwind CSS。\n- 支持响应式布局，适配移动端。\n- 支持暗色模式。\n\nAPI 响应示例：\n[在此粘贴 OpenAPI 定义或 JSON 示例]\n为什么这样分工？\n\nCodex 负责逻辑和数据（后端、API 设计）。\n\nTrae 负责视觉和交互（前端 UI），两者通过 API 文档解耦，互不干扰。\n\n🔧 第三阶段：终端辅助与 Git 管理（OpenCode）\nOpenCode 在你的工作流中扮演“特种兵”角色，处理那些小而确定、但频繁打断主流程的任务。\n\n场景 3.1：快速修复已知 Bug\n不进入 Codex 重上下文，直接在终端执行：\n\nbash\nopencode \"修复 src/utils/date.js 第 42 行的空指针异常，只输出 patch 文件内容，不要直接修改文件\"\n拿到 patch 后，你可以先审阅再手动应用，安全性极高。\n\n场景 3.2：智能生成 Git 提交信息\n当暂存区准备好后，在终端执行：\n\nbash\nopencode \"分析当前 git diff --staged 的变更内容，生成一条符合 Conventional Commits 规范的提交信息，包含清晰的 body 说明\"\nOpenCode 会输出类似：\n\ntext\nfeat(points): add user points expiration logic\n\n- Add `expires_at` column to `user_points` table\n- Implement daily cleanup job for expired points\n- Update points balance query to filter out expired records\n场景 3.3：用平板安全浏览和审查代码\n当你通过平板远程 SSH 到笔记本时，使用 --dry-run 模式保证只读：\n\nbash\nopencode --dry-run \"分析当前分支中与积分相关的所有代码文件，列出潜在的性能瓶颈和未处理的边界条件。仅输出报告，不修改任何文件。\"\n这相当于在平板上用 AI 做了一次免费的 Code Review。\n\n🔄 第四阶段：中后期维护 —— 应对“需求变更”和“文档混乱”\n这是整个工作流中最体现工程化水平的一环。\n\n场景 4.1：开发中途，需求发生变更（文档联动更新）\n问题： 代码写了一半，突然觉得积分过期规则不合理，想从“年底清零”改为“获得后 365 天过期”。\n\n错误做法： 手动去改代码，然后忘了改 PRD 和设计文档。\n\n正确做法： 在 Codex 中下达联动修改指令。\n\nmarkdown\n我决定修改积分过期规则：从“每年12月31日清零”改为“自获得之日起365天后过期”。\n\n请**同步修改**以下文件：\n1. `/docs/prd/prd-积分系统.md` 中的业务规则部分。\n2. `/docs/design/design-积分系统.md` 中的数据模型和 API 说明。\n3. 后端代码中对应的逻辑实现（包括数据库迁移、业务层、测试用例）。\n原理： Codex（GPT-5.4）具备强大的上下文理解能力，能一次性在多个文件间保持修改的一致性。你只需要做“决策”，AI 负责“扩散修改”。\n\n场景 4.2：已经改乱了，代码和文档严重脱节（逆向同步）\n问题： 你手快改了一堆代码，并且已经提交了好几次。现在 /docs/design.md 完全是一张废纸，继续开发时 AI 频繁产生幻觉。\n\n解决方法： 用 Codex 做逆向工程。\n\n步骤：\n\n将当前核心代码文件的内容（如 models/user_points.go、api/points_handler.go）复制或通过 @ 引用给 Codex。\n\n输入以下指令：\n\nmarkdown\n请忽略旧的 `/docs/design/design-积分系统.md`。\n\n根据以下**当前正在运行的代码**，反向生成一份准确的技术设计文档。要求覆盖现有实现的所有数据表、API 接口和核心逻辑。输出为新文件 `/docs/design/design-积分系统-v2.md`。\n效果： 几分钟内，你就能获得一份与代码 100% 吻合的新文档。拿它继续喂给 AI，幻觉问题立刻消失。\n\n📱 第五阶段：移动办公增强（平板场景）\n你提到常用平板远程操控笔记本。以下技巧可将平板从“远程屏幕”升级为“AI 指挥中心”。\n\n技巧 5.1：建立 tmux 后台任务\n通过 SSH 连上笔记本后，启动一个 tmux 会话，并在其中运行 Codex 的后台任务。\n\nbash\ntmux new -s ai-worker\ncodex --task \"重构 src/services/ 目录下的模块，提高内聚性，完成后输出总结报告\"\n# 按 Ctrl+B D 分离会话，关闭平板，任务在笔记本后台持续运行\n下次连上后执行 tmux attach -t ai-worker 即可查看结果。\n\n技巧 5.2：使用 OpenCode Web 模式（可选）\n在笔记本终端运行：\n\nbash\nopencode --server --port 8899\n然后在平板浏览器访问 http://笔记本IP:8899，你会得到一个 OpenCode 的 Web UI。你可以直接在这个界面中下达终端指令、查看 Git 状态，触屏操作比打字更便捷。\n\n💎 总结：这套工作流解决了什么问题？\n规划质量：用最强模型（GPT-5.4）统一产出 PRD 和设计文档，避免多工具切换导致的思路中断和文档不一致。\n\n执行效率：Codex 负责重逻辑，Trae 负责前端 UI，OpenCode 负责轻量杂务，各司其职，互不阻塞。\n\n中后期熵增：通过“联动修改指令”和“逆向文档同步”两个杀手锏，彻底解决“代码和文档对不上 → AI 幻觉 → 效率崩塌”的恶性循环。\n\n移动办公：利用 tmux 和 OpenCode 的远程特性，将平板变成 AI 任务调度中心。\n\n📝 面试时的一句话总结\n“我搭建了一套以 Codex 为核心的 AI 原生开发流水线。PRD、技术设计和全栈编码均由 GPT-5.4 统一驱动，确保从需求到实现的强一致性。同时，我制定了一套文档与代码的联动更新规范，并通过 Trae 和 OpenCode 进行专项任务分流，将 AI 幻觉率控制在极低水平。整个环境支持远程接入，在平板上即可完成代码审查与任务调度。”\n\n后续优化方向建议：\n\n将常用的联动修改指令保存为 Codex 的 Skill，一键触发。\n\n尝试在项目中加入 CODEBUDDY.md 或 .cursorrules 文件，将你的开发规范（如提交信息格式、文档存放路径）固化下来，让 AI 自动遵守。\n\n如果你在使用过程中发现了新的技巧或痛点，欢迎随时补充到这篇文档里。','/uploads/1776251801801137200.png','published','',4,'2026-04-15 19:18:47.497','2026-04-15 19:04:34.690','2026-04-15 19:18:47.498',3,4),(5,'Test Draft Article','Testing','This is a test draft','','published','',1,'2026-06-01 12:38:25.655','2026-06-01 12:38:25.655','2026-06-01 12:38:25.660',NULL,1),(7,'Updated Title','Updated','Updated','','published','',1,'2026-06-01 12:38:49.782','2026-06-01 12:38:49.782','2026-06-01 12:38:50.065',NULL,2),(8,'QA��������','','# ������������\n\n����һƪQA�������µ����ݡ�\n\n## �����½�\n\n- Ҫ��1\n- Ҫ��2\n- Ҫ��3','','published','',1,'2026-06-01 14:07:02.665','2026-06-01 14:07:02.665','2026-06-01 14:07:02.670',1,1),(9,'QA��������','','# ������������\n\n����һƪQA�������µ����ݡ�\n\n## �����½�\n\n- Ҫ��1\n- Ҫ��2\n- Ҫ��3','','published','',1,'2026-06-01 14:07:05.958','2026-06-01 14:07:05.958','2026-06-01 14:07:05.962',1,2),(13,'Review Test Article','','# Please review','','published','',6,'2026-06-01 14:08:37.101','2026-06-01 14:08:36.543','2026-06-01 14:08:37.101',1,1);
/*!40000 ALTER TABLE `articles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `categories`
--

DROP TABLE IF EXISTS `categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `slug` varchar(60) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_categories_name` (`name`),
  UNIQUE KEY `idx_categories_slug` (`slug`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `categories`
--

LOCK TABLES `categories` WRITE;
/*!40000 ALTER TABLE `categories` DISABLE KEYS */;
INSERT INTO `categories` VALUES (1,'General','default','2026-04-02 13:10:27.631','2026-04-08 12:15:01.025'),(2,'Tech Notes','tech-notes','2026-04-02 13:10:27.653','2026-04-08 12:15:01.036'),(3,'Product Design','product-design','2026-04-02 13:10:27.657','2026-04-08 12:15:01.041'),(4,'dwadw','dwadw','2026-04-02 13:39:53.537','2026-04-02 13:39:53.537'),(5,'Test Cat','test-cat','2026-06-01 12:38:50.514','2026-06-01 12:38:50.514'),(6,'QA���Է���','qa���է���','2026-06-01 14:09:49.945','2026-06-01 14:09:49.945'),(7,'QA Category','qa-category','2026-06-01 14:10:23.956','2026-06-01 14:10:23.956');
/*!40000 ALTER TABLE `categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(500) NOT NULL,
  `article_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `parent_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_comments_article_id` (`article_id`),
  KEY `idx_comments_user_id` (`user_id`),
  KEY `idx_comments_parent_id` (`parent_id`),
  CONSTRAINT `fk_articles_comments` FOREIGN KEY (`article_id`) REFERENCES `articles` (`id`),
  CONSTRAINT `fk_comments_replies` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`),
  CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (6,'Test comment',7,1,'2026-06-01 12:38:50.418','2026-06-01 12:38:50.418',NULL);
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `moderation_hits`
--

DROP TABLE IF EXISTS `moderation_hits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `moderation_hits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `scene` varchar(40) NOT NULL,
  `field` varchar(80) NOT NULL,
  `matched_word` varchar(120) NOT NULL,
  `snippet` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_moderation_hits_user_id` (`user_id`),
  KEY `idx_moderation_hits_scene` (`scene`),
  CONSTRAINT `fk_moderation_hits_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `moderation_hits`
--

LOCK TABLES `moderation_hits` WRITE;
/*!40000 ALTER TABLE `moderation_hits` DISABLE KEYS */;
/*!40000 ALTER TABLE `moderation_hits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `title` varchar(120) NOT NULL,
  `content` varchar(255) NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `type` varchar(40) NOT NULL DEFAULT 'system',
  `action_url` varchar(255) DEFAULT NULL,
  `payload` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_user_id` (`user_id`),
  KEY `idx_notifications_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifications`
--

LOCK TABLES `notifications` WRITE;
/*!40000 ALTER TABLE `notifications` DISABLE KEYS */;
INSERT INTO `notifications` VALUES (1,3,'收到新评论','你的文章《dwada》收到了新的评论。',0,'2026-04-03 16:07:41.992','system',NULL,NULL),(2,3,'收到新评论','你的文章《dwada》收到了新的评论。',0,'2026-04-03 16:07:42.002','system',NULL,NULL),(3,3,'收到新评论','你的文章《dwada》收到了新的评论。',0,'2026-04-03 16:08:52.760','system',NULL,NULL),(4,3,'收到新评论','你的文章《dwada》收到了新的评论。',0,'2026-04-03 16:08:52.774','system',NULL,NULL),(5,4,'审核结果通知','你的文章《以 Codex 为中心的 AI 全栈开发工作流：从 PRD 到代码联动》已审核通过并发布。',0,'2026-04-15 19:18:47.515','system',NULL,NULL),(6,6,'审核结果通知','你的文章《Review Test Article》已审核通过并发布。',0,'2026-06-01 14:08:37.109','article_review','/my-articles','{\"action\":\"approve\",\"articleId\":13}'),(7,6,'审核结果通知','你的文章《Reject Test》未通过审核。',0,'2026-06-01 14:08:37.523','article_review','/my-articles','{\"action\":\"reject\",\"articleId\":14}'),(8,1,'收到新评论','你的文章收到了新的评论：QA��������',0,'2026-06-01 14:08:49.079','article_comment','/article/9','{\"articleId\":9,\"commentId\":7}'),(9,1,'收到新评论','你的文章收到了新的评论：QA��������',0,'2026-06-01 14:08:54.813','article_comment','/article/9','{\"articleId\":9,\"commentId\":8}'),(10,1,'收到新评论','你的文章收到了新的评论：QA��������',0,'2026-06-01 14:09:35.294','article_comment','/article/9','{\"articleId\":9,\"commentId\":9}'),(11,6,'审核结果通知','你的项目《User Project Draft》已审核通过并发布。',0,'2026-06-01 14:10:23.880','project_review','/my-projects','{\"action\":\"approve\",\"projectId\":9}');
/*!40000 ALTER TABLE `notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `operation_logs`
--

DROP TABLE IF EXISTS `operation_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `operator_id` bigint unsigned NOT NULL,
  `action` varchar(80) NOT NULL,
  `target_type` varchar(40) NOT NULL,
  `target_id` bigint unsigned NOT NULL,
  `description` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_operation_logs_operator_id` (`operator_id`),
  CONSTRAINT `fk_operation_logs_operator` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=233 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `operation_logs`
--

LOCK TABLES `operation_logs` WRITE;
/*!40000 ALTER TABLE `operation_logs` DISABLE KEYS */;
INSERT INTO `operation_logs` VALUES (1,1,'delete_upload','upload',0,'删除图片资源：admin-delete-test.txt','2026-04-02 19:17:44.924'),(2,1,'delete_comment','comment',2,'删除评论','2026-04-03 16:08:05.491'),(3,1,'delete_comment','comment',4,'删除评论','2026-04-03 16:08:52.783'),(4,1,'review_article','article',4,'你的文章《以 Codex 为中心的 AI 全栈开发工作流：从 PRD 到代码联动》已审核通过并发布。','2026-04-15 19:18:47.518'),(5,1,'update_user_role','user',4,'修改用户角色为 admin','2026-04-15 19:18:55.561'),(6,1,'delete_article','article',3,'删除文章','2026-04-15 19:23:28.027'),(7,1,'delete_article','article',2,'删除文章','2026-04-15 19:23:31.913'),(8,1,'delete_article','article',1,'删除文章','2026-04-15 19:23:35.250'),(9,1,'delete_upload','upload',0,'删除图片资源：1776251771473396200.png','2026-04-15 19:23:58.274'),(10,1,'delete_upload','upload',0,'删除图片资源：1776251211541250800.png','2026-04-15 19:24:00.516'),(11,1,'delete_upload','upload',0,'删除图片资源：1775108312286743800.jpg','2026-04-15 19:24:02.896'),(12,1,'update_moderation_settings','system_setting',0,'更新自动封禁阈值','2026-06-01 12:40:35.759'),(13,1,'delete_article','article',10,'删除文章','2026-06-01 14:07:49.504'),(14,1,'review_article','article',13,'你的文章《Review Test Article》已审核通过并发布。','2026-06-01 14:08:37.111'),(15,1,'review_article','article',14,'你的文章《Reject Test》未通过审核。','2026-06-01 14:08:37.524'),(16,1,'delete_comment','comment',8,'删除评论','2026-06-01 14:09:00.748'),(17,1,'delete_comment','comment',9,'删除评论','2026-06-01 14:09:35.542'),(18,1,'delete_project','project',7,'删除项目','2026-06-01 14:10:08.648'),(19,1,'review_project','project',9,'你的项目《User Project Draft》已审核通过并发布。','2026-06-01 14:10:23.882'),(20,1,'update_moderation_settings','system_setting',0,'更新自动封禁阈值','2026-06-01 14:10:41.761'),(21,1,'update_user_status','user',6,'修改用户状态为 banned','2026-06-01 14:10:41.903'),(22,1,'update_user_status','user',6,'修改用户状态为 active','2026-06-01 14:10:42.035'),(23,1,'delete_article','article',16,'删除文章','2026-06-01 14:11:31.189'),(24,1,'delete_article','article',15,'删除文章','2026-06-01 14:11:31.235'),(25,1,'delete_article','article',14,'删除文章','2026-06-01 14:11:31.279'),(26,1,'delete_article','article',12,'删除文章','2026-06-01 14:11:31.346'),(27,1,'delete_project','project',8,'删除项目','2026-06-01 14:11:31.399'),(28,1,'delete_project','project',9,'删除项目','2026-06-01 14:11:31.454'),(29,1,'delete_comment','comment',7,'删除评论','2026-06-01 14:11:31.496'),(30,1,'delete_article','article',18,'删除文章','2026-06-02 00:06:27.544'),(31,1,'delete_project','project',10,'删除项目','2026-06-02 00:06:27.590'),(32,1,'delete_article','article',17,'删除文章','2026-06-04 13:29:33.164'),(33,1,'delete_article','article',19,'删除文章','2026-06-04 13:29:37.239'),(34,1,'delete_article','article',11,'删除文章','2026-06-04 13:29:39.844'),(35,1,'delete_article','article',6,'删除文章','2026-06-04 13:29:44.764'),(36,1,'create_sensitive_word','sensitive_word',1,'新增敏感词：六四事件','2026-06-04 13:34:40.048'),(37,1,'create_sensitive_word','sensitive_word',2,'新增敏感词：天安门母亲','2026-06-04 13:34:40.057'),(38,1,'create_sensitive_word','sensitive_word',3,'新增敏感词：极端宗教','2026-06-04 13:34:40.062'),(39,1,'create_sensitive_word','sensitive_word',4,'新增敏感词：宗教极端','2026-06-04 13:34:40.068'),(40,1,'create_sensitive_word','sensitive_word',5,'新增敏感词：民族分裂','2026-06-04 13:34:40.073'),(41,1,'create_sensitive_word','sensitive_word',6,'新增敏感词：宗教迫害','2026-06-04 13:34:40.077'),(42,1,'create_sensitive_word','sensitive_word',7,'新增敏感词：非法组织','2026-06-04 13:34:40.083'),(43,1,'create_sensitive_word','sensitive_word',8,'新增敏感词：非法集会','2026-06-04 13:34:40.087'),(44,1,'create_sensitive_word','sensitive_word',9,'新增敏感词：非法游行','2026-06-04 13:34:40.094'),(45,1,'create_sensitive_word','sensitive_word',10,'新增敏感词：非法示威','2026-06-04 13:34:40.099'),(46,1,'create_sensitive_word','sensitive_word',11,'新增敏感词：非法聚集','2026-06-04 13:34:40.104'),(47,1,'create_sensitive_word','sensitive_word',12,'新增敏感词：静坐','2026-06-04 13:34:40.109'),(48,1,'create_sensitive_word','sensitive_word',13,'新增敏感词：拉横幅','2026-06-04 13:34:40.114'),(49,1,'create_sensitive_word','sensitive_word',14,'新增敏感词：越级上访','2026-06-04 13:34:40.118'),(50,1,'create_sensitive_word','sensitive_word',15,'新增敏感词：被消失','2026-06-04 13:34:40.123'),(51,1,'create_sensitive_word','sensitive_word',16,'新增敏感词：被自杀','2026-06-04 13:34:40.128'),(52,1,'create_sensitive_word','sensitive_word',17,'新增敏感词：恶意讨薪','2026-06-04 13:34:40.133'),(53,1,'create_sensitive_word','sensitive_word',18,'新增敏感词：恐怖主义','2026-06-04 13:34:40.137'),(54,1,'create_sensitive_word','sensitive_word',19,'新增敏感词：极端主义','2026-06-04 13:34:40.143'),(55,1,'create_sensitive_word','sensitive_word',20,'新增敏感词：自焚','2026-06-04 13:34:40.149'),(56,1,'create_sensitive_word','sensitive_word',21,'新增敏感词：割腕','2026-06-04 13:34:40.154'),(57,1,'create_sensitive_word','sensitive_word',22,'新增敏感词：服毒','2026-06-04 13:34:40.160'),(58,1,'create_sensitive_word','sensitive_word',23,'新增敏感词：安眠药','2026-06-04 13:34:40.165'),(59,1,'create_sensitive_word','sensitive_word',24,'新增敏感词：百草枯','2026-06-04 13:34:40.170'),(60,1,'create_sensitive_word','sensitive_word',25,'新增敏感词：氰化物','2026-06-04 13:34:40.175'),(61,1,'create_sensitive_word','sensitive_word',26,'新增敏感词：杀人灭口','2026-06-04 13:34:40.181'),(62,1,'create_sensitive_word','sensitive_word',27,'新增敏感词：绑架','2026-06-04 13:34:40.185'),(63,1,'create_sensitive_word','sensitive_word',28,'新增敏感词：分尸','2026-06-04 13:34:40.191'),(64,1,'create_sensitive_word','sensitive_word',29,'新增敏感词：碎尸','2026-06-04 13:34:40.195'),(65,1,'create_sensitive_word','sensitive_word',30,'新增敏感词：肢解','2026-06-04 13:34:40.200'),(66,1,'create_sensitive_word','sensitive_word',31,'新增敏感词：斩首','2026-06-04 13:34:40.205'),(67,1,'create_sensitive_word','sensitive_word',32,'新增敏感词：暴恐音视频','2026-06-04 13:34:40.210'),(68,1,'create_sensitive_word','sensitive_word',33,'新增敏感词：血腥暴力','2026-06-04 13:34:40.214'),(69,1,'create_sensitive_word','sensitive_word',34,'新增敏感词：迷药','2026-06-04 13:34:40.220'),(70,1,'create_sensitive_word','sensitive_word',35,'新增敏感词：迷魂药','2026-06-04 13:34:40.231'),(71,1,'create_sensitive_word','sensitive_word',36,'新增敏感词：催情药','2026-06-04 13:34:40.237'),(72,1,'create_sensitive_word','sensitive_word',37,'新增敏感词：听话水','2026-06-04 13:34:40.243'),(73,1,'create_sensitive_word','sensitive_word',38,'新增敏感词：神仙水','2026-06-04 13:34:40.249'),(74,1,'create_sensitive_word','sensitive_word',39,'新增敏感词：摇头丸','2026-06-04 13:34:40.254'),(75,1,'create_sensitive_word','sensitive_word',40,'新增敏感词：麻古','2026-06-04 13:34:40.261'),(76,1,'create_sensitive_word','sensitive_word',41,'新增敏感词：止咳水','2026-06-04 13:34:40.268'),(77,1,'create_sensitive_word','sensitive_word',42,'新增敏感词：蓝精灵','2026-06-04 13:34:40.274'),(78,1,'create_sensitive_word','sensitive_word',43,'新增敏感词：邮票毒品','2026-06-04 13:34:40.279'),(79,1,'create_sensitive_word','sensitive_word',44,'新增敏感词：电子烟毒品','2026-06-04 13:34:40.284'),(80,1,'create_sensitive_word','sensitive_word',45,'新增敏感词：制毒','2026-06-04 13:34:40.289'),(81,1,'create_sensitive_word','sensitive_word',46,'新增敏感词：贩毒','2026-06-04 13:34:40.297'),(82,1,'create_sensitive_word','sensitive_word',47,'新增敏感词：仿真枪','2026-06-04 13:34:40.302'),(83,1,'create_sensitive_word','sensitive_word',48,'新增敏感词：气枪','2026-06-04 13:34:40.307'),(84,1,'create_sensitive_word','sensitive_word',49,'新增敏感词：猎枪','2026-06-04 13:34:40.314'),(85,1,'create_sensitive_word','sensitive_word',50,'新增敏感词：钢珠枪','2026-06-04 13:34:40.319'),(86,1,'create_sensitive_word','sensitive_word',51,'新增敏感词：电击棒','2026-06-04 13:34:40.324'),(87,1,'create_sensitive_word','sensitive_word',52,'新增敏感词：弹簧刀','2026-06-04 13:34:40.329'),(88,1,'create_sensitive_word','sensitive_word',53,'新增敏感词：匕首','2026-06-04 13:34:40.334'),(89,1,'create_sensitive_word','sensitive_word',54,'新增敏感词：卖肾','2026-06-04 13:34:40.339'),(90,1,'create_sensitive_word','sensitive_word',55,'新增敏感词：器官买卖','2026-06-04 13:34:40.343'),(91,1,'create_sensitive_word','sensitive_word',56,'新增敏感词：代孕','2026-06-04 13:34:40.349'),(92,1,'create_sensitive_word','sensitive_word',57,'新增敏感词：卵子买卖','2026-06-04 13:34:40.354'),(93,1,'create_sensitive_word','sensitive_word',58,'新增敏感词：三级片','2026-06-04 13:35:03.306'),(94,1,'create_sensitive_word','sensitive_word',59,'新增敏感词：成人片','2026-06-04 13:35:03.314'),(95,1,'create_sensitive_word','sensitive_word',60,'新增敏感词：毛片','2026-06-04 13:35:03.320'),(96,1,'create_sensitive_word','sensitive_word',61,'新增敏感词：无码','2026-06-04 13:35:03.325'),(97,1,'create_sensitive_word','sensitive_word',62,'新增敏感词：成人电影','2026-06-04 13:35:03.330'),(98,1,'create_sensitive_word','sensitive_word',63,'新增敏感词：情色','2026-06-04 13:35:03.334'),(99,1,'create_sensitive_word','sensitive_word',64,'新增敏感词：裸模','2026-06-04 13:35:03.339'),(100,1,'create_sensitive_word','sensitive_word',65,'新增敏感词：一夜情','2026-06-04 13:35:03.345'),(101,1,'create_sensitive_word','sensitive_word',66,'新增敏感词：同城约炮','2026-06-04 13:35:03.350'),(102,1,'create_sensitive_word','sensitive_word',67,'新增敏感词：淫秽','2026-06-04 13:35:03.354'),(103,1,'create_sensitive_word','sensitive_word',68,'新增敏感词：淫乱','2026-06-04 13:35:03.359'),(104,1,'create_sensitive_word','sensitive_word',69,'新增敏感词：淫荡','2026-06-04 13:35:03.364'),(105,1,'create_sensitive_word','sensitive_word',70,'新增敏感词：性交易','2026-06-04 13:35:03.369'),(106,1,'create_sensitive_word','sensitive_word',71,'新增敏感词：包养','2026-06-04 13:35:03.373'),(107,1,'create_sensitive_word','sensitive_word',72,'新增敏感词：失足少女','2026-06-04 13:35:03.377'),(108,1,'create_sensitive_word','sensitive_word',73,'新增敏感词：陪睡','2026-06-04 13:35:03.382'),(109,1,'create_sensitive_word','sensitive_word',74,'新增敏感词：出台','2026-06-04 13:35:03.387'),(110,1,'create_sensitive_word','sensitive_word',75,'新增敏感词：坐台','2026-06-04 13:35:03.392'),(111,1,'create_sensitive_word','sensitive_word',76,'新增敏感词：裸体','2026-06-04 13:35:03.397'),(112,1,'create_sensitive_word','sensitive_word',77,'新增敏感词：裸照','2026-06-04 13:35:03.401'),(113,1,'create_sensitive_word','sensitive_word',78,'新增敏感词：偷拍','2026-06-04 13:35:03.407'),(114,1,'create_sensitive_word','sensitive_word',79,'新增敏感词：走光','2026-06-04 13:35:03.413'),(115,1,'create_sensitive_word','sensitive_word',80,'新增敏感词：艳照','2026-06-04 13:35:03.417'),(116,1,'create_sensitive_word','sensitive_word',81,'新增敏感词：迷奸','2026-06-04 13:35:03.423'),(117,1,'create_sensitive_word','sensitive_word',82,'新增敏感词：兽交','2026-06-04 13:35:03.428'),(118,1,'create_sensitive_word','sensitive_word',83,'新增敏感词：乱伦','2026-06-04 13:35:03.433'),(119,1,'create_sensitive_word','sensitive_word',84,'新增敏感词：幼女','2026-06-04 13:35:03.438'),(120,1,'create_sensitive_word','sensitive_word',85,'新增敏感词：少女色情','2026-06-04 13:35:03.443'),(121,1,'create_sensitive_word','sensitive_word',86,'新增敏感词：萝莉','2026-06-04 13:35:03.449'),(122,1,'create_sensitive_word','sensitive_word',87,'新增敏感词：性奴','2026-06-04 13:35:03.455'),(123,1,'create_sensitive_word','sensitive_word',88,'新增敏感词：sm','2026-06-04 13:35:03.460'),(124,1,'create_sensitive_word','sensitive_word',89,'新增敏感词：nsfw','2026-06-04 13:35:03.464'),(125,1,'create_sensitive_word','sensitive_word',90,'新增敏感词：onlyfans','2026-06-04 13:35:03.470'),(126,1,'create_sensitive_word','sensitive_word',91,'新增敏感词：时时彩','2026-06-04 13:35:03.474'),(127,1,'create_sensitive_word','sensitive_word',92,'新增敏感词：六合彩','2026-06-04 13:35:03.479'),(128,1,'create_sensitive_word','sensitive_word',93,'新增敏感词：老虎机','2026-06-04 13:35:03.484'),(129,1,'create_sensitive_word','sensitive_word',94,'新增敏感词：百家乐','2026-06-04 13:35:03.489'),(130,1,'create_sensitive_word','sensitive_word',95,'新增敏感词：德州扑克','2026-06-04 13:35:03.494'),(131,1,'create_sensitive_word','sensitive_word',96,'新增敏感词：牌九','2026-06-04 13:35:03.499'),(132,1,'create_sensitive_word','sensitive_word',97,'新增敏感词：梭哈','2026-06-04 13:35:03.503'),(133,1,'create_sensitive_word','sensitive_word',98,'新增敏感词：赌球','2026-06-04 13:35:03.508'),(134,1,'create_sensitive_word','sensitive_word',99,'新增敏感词：赌马','2026-06-04 13:35:03.513'),(135,1,'create_sensitive_word','sensitive_word',100,'新增敏感词：赌船','2026-06-04 13:35:03.519'),(136,1,'create_sensitive_word','sensitive_word',101,'新增敏感词：网络赌场','2026-06-04 13:35:03.524'),(137,1,'create_sensitive_word','sensitive_word',102,'新增敏感词：棋牌赌博','2026-06-04 13:35:03.529'),(138,1,'create_sensitive_word','sensitive_word',103,'新增敏感词：捕鱼游戏','2026-06-04 13:35:03.537'),(139,1,'create_sensitive_word','sensitive_word',104,'新增敏感词：电子游艺','2026-06-04 13:35:03.543'),(140,1,'create_sensitive_word','sensitive_word',105,'新增敏感词：体育投注','2026-06-04 13:35:03.549'),(141,1,'create_sensitive_word','sensitive_word',106,'新增敏感词：电子赌博','2026-06-04 13:35:03.555'),(142,1,'create_sensitive_word','sensitive_word',107,'新增敏感词：网络赌博','2026-06-04 13:35:03.562'),(143,1,'create_sensitive_word','sensitive_word',108,'新增敏感词：投注平台','2026-06-04 13:35:03.568'),(144,1,'create_sensitive_word','sensitive_word',109,'新增敏感词：现金棋牌','2026-06-04 13:35:03.573'),(145,1,'create_sensitive_word','sensitive_word',110,'新增敏感词：信誉平台','2026-06-04 13:35:03.578'),(146,1,'create_sensitive_word','sensitive_word',111,'新增敏感词：真人娱乐','2026-06-04 13:35:03.583'),(147,1,'create_sensitive_word','sensitive_word',112,'新增敏感词：澳门赌场','2026-06-04 13:35:03.590'),(148,1,'create_sensitive_word','sensitive_word',113,'新增敏感词：葡京','2026-06-04 13:35:03.595'),(149,1,'create_sensitive_word','sensitive_word',114,'新增敏感词：网络传销','2026-06-04 13:35:03.600'),(150,1,'create_sensitive_word','sensitive_word',115,'新增敏感词：资金盘','2026-06-04 13:35:03.605'),(151,1,'create_sensitive_word','sensitive_word',116,'新增敏感词：杀猪盘','2026-06-04 13:35:03.610'),(152,1,'create_sensitive_word','sensitive_word',117,'新增敏感词：庞氏骗局','2026-06-04 13:35:03.616'),(153,1,'create_sensitive_word','sensitive_word',118,'新增敏感词：拉人头','2026-06-04 13:35:03.622'),(154,1,'create_sensitive_word','sensitive_word',119,'新增敏感词：零风险高收益','2026-06-04 13:35:03.628'),(155,1,'create_sensitive_word','sensitive_word',120,'新增敏感词：稳赚不赔','2026-06-04 13:35:03.633'),(156,1,'create_sensitive_word','sensitive_word',121,'新增敏感词：日赚','2026-06-04 13:35:03.638'),(157,1,'create_sensitive_word','sensitive_word',122,'新增敏感词：月入过万','2026-06-04 13:35:03.643'),(158,1,'create_sensitive_word','sensitive_word',123,'新增敏感词：挂机赚钱','2026-06-04 13:35:03.649'),(159,1,'create_sensitive_word','sensitive_word',124,'新增敏感词：点赞赚钱','2026-06-04 13:35:03.654'),(160,1,'create_sensitive_word','sensitive_word',125,'新增敏感词：虚拟币诈骗','2026-06-04 13:35:03.659'),(161,1,'create_sensitive_word','sensitive_word',126,'新增敏感词：区块链骗局','2026-06-04 13:35:03.665'),(162,1,'create_sensitive_word','sensitive_word',127,'新增敏感词：投资返利','2026-06-04 13:35:03.671'),(163,1,'create_sensitive_word','sensitive_word',128,'新增敏感词：推荐股票','2026-06-04 13:35:03.676'),(164,1,'create_sensitive_word','sensitive_word',129,'新增敏感词：内幕消息','2026-06-04 13:35:03.682'),(165,1,'create_sensitive_word','sensitive_word',130,'新增敏感词：割韭菜','2026-06-04 13:35:03.687'),(166,1,'create_sensitive_word','sensitive_word',131,'新增敏感词：薅羊毛','2026-06-04 13:35:03.692'),(167,1,'create_sensitive_word','sensitive_word',132,'新增敏感词：接码平台','2026-06-04 13:35:03.699'),(168,1,'create_sensitive_word','sensitive_word',133,'新增敏感词：猫池','2026-06-04 13:35:03.704'),(169,1,'create_sensitive_word','sensitive_word',134,'新增敏感词：伪基站','2026-06-04 13:35:03.709'),(170,1,'create_sensitive_word','sensitive_word',135,'新增敏感词：短信轰炸','2026-06-04 13:35:03.714'),(171,1,'create_sensitive_word','sensitive_word',136,'新增敏感词：电话轰炸','2026-06-04 13:35:03.720'),(172,1,'create_sensitive_word','sensitive_word',137,'新增敏感词：种族歧视','2026-06-04 13:36:04.075'),(173,1,'create_sensitive_word','sensitive_word',138,'新增敏感词：地域黑','2026-06-04 13:36:04.083'),(174,1,'create_sensitive_word','sensitive_word',139,'新增敏感词：地域歧视','2026-06-04 13:36:04.090'),(175,1,'create_sensitive_word','sensitive_word',140,'新增敏感词：歧视女性','2026-06-04 13:36:04.096'),(176,1,'create_sensitive_word','sensitive_word',141,'新增敏感词：歧视残疾人','2026-06-04 13:36:04.102'),(177,1,'create_sensitive_word','sensitive_word',142,'新增敏感词：人身攻击','2026-06-04 13:36:04.107'),(178,1,'create_sensitive_word','sensitive_word',143,'新增敏感词：网络暴力','2026-06-04 13:36:04.112'),(179,1,'create_sensitive_word','sensitive_word',144,'新增敏感词：人肉搜索','2026-06-04 13:36:04.119'),(180,1,'create_sensitive_word','sensitive_word',145,'新增敏感词：开盒','2026-06-04 13:36:04.125'),(181,1,'create_sensitive_word','sensitive_word',146,'新增敏感词：死亡威胁','2026-06-04 13:36:04.130'),(182,1,'create_sensitive_word','sensitive_word',147,'新增敏感词：去死','2026-06-04 13:36:04.136'),(183,1,'create_sensitive_word','sensitive_word',148,'新增敏感词：全家死','2026-06-04 13:36:04.142'),(184,1,'create_sensitive_word','sensitive_word',149,'新增敏感词：畜生','2026-06-04 13:36:04.148'),(185,1,'create_sensitive_word','sensitive_word',150,'新增敏感词：娘炮','2026-06-04 13:36:04.153'),(186,1,'create_sensitive_word','sensitive_word',151,'新增敏感词：脱库','2026-06-04 13:36:04.159'),(187,1,'create_sensitive_word','sensitive_word',152,'新增敏感词：撞库','2026-06-04 13:36:04.164'),(188,1,'create_sensitive_word','sensitive_word',153,'新增敏感词：删库跑路','2026-06-04 13:36:04.169'),(189,1,'create_sensitive_word','sensitive_word',154,'新增敏感词：webshell','2026-06-04 13:36:04.175'),(190,1,'create_sensitive_word','sensitive_word',155,'新增敏感词：卷款潜逃','2026-06-04 13:36:04.180'),(191,1,'create_sensitive_word','sensitive_word',156,'新增敏感词：非法经营','2026-06-04 13:36:04.186'),(192,1,'create_sensitive_word','sensitive_word',157,'新增敏感词：无证经营','2026-06-04 13:36:04.191'),(193,1,'create_sensitive_word','sensitive_word',158,'新增敏感词：偷税漏税','2026-06-04 13:36:04.197'),(194,1,'create_sensitive_word','sensitive_word',159,'新增敏感词：走私','2026-06-04 13:36:04.202'),(195,1,'create_sensitive_word','sensitive_word',160,'新增敏感词：地下钱庄','2026-06-04 13:36:04.207'),(196,1,'create_sensitive_word','sensitive_word',161,'新增敏感词：高利贷','2026-06-04 13:36:04.212'),(197,1,'create_sensitive_word','sensitive_word',162,'新增敏感词：套路贷','2026-06-04 13:36:04.218'),(198,1,'create_sensitive_word','sensitive_word',163,'新增敏感词：校园贷','2026-06-04 13:36:04.224'),(199,1,'create_sensitive_word','sensitive_word',164,'新增敏感词：裸贷','2026-06-04 13:36:04.229'),(200,1,'create_sensitive_word','sensitive_word',165,'新增敏感词：暴力催收','2026-06-04 13:36:04.234'),(201,1,'create_sensitive_word','sensitive_word',166,'新增敏感词：盗版资源','2026-06-04 13:36:04.240'),(202,1,'create_sensitive_word','sensitive_word',167,'新增敏感词：注册机','2026-06-04 13:36:04.245'),(203,1,'create_sensitive_word','sensitive_word',168,'新增敏感词：游戏外挂','2026-06-04 13:36:04.251'),(204,1,'create_sensitive_word','sensitive_word',169,'新增敏感词：私服','2026-06-04 13:36:04.256'),(205,1,'create_sensitive_word','sensitive_word',170,'新增敏感词：电信诈骗','2026-06-04 13:36:04.263'),(206,1,'create_sensitive_word','sensitive_word',171,'新增敏感词：帮信罪','2026-06-04 13:36:04.269'),(207,1,'create_sensitive_word','sensitive_word',172,'新增敏感词：跨境赌博','2026-06-04 13:36:04.275'),(208,1,'create_sensitive_word','sensitive_word',173,'新增敏感词：代写论文','2026-06-04 13:36:04.280'),(209,1,'create_sensitive_word','sensitive_word',174,'新增敏感词：代发论文','2026-06-04 13:36:04.285'),(210,1,'create_sensitive_word','sensitive_word',175,'新增敏感词：论文代写','2026-06-04 13:36:04.290'),(211,1,'create_sensitive_word','sensitive_word',176,'新增敏感词：刷粉','2026-06-04 13:36:04.295'),(212,1,'create_sensitive_word','sensitive_word',177,'新增敏感词：刷赞','2026-06-04 13:36:04.301'),(213,1,'create_sensitive_word','sensitive_word',178,'新增敏感词：刷屏','2026-06-04 13:36:04.307'),(214,1,'create_sensitive_word','sensitive_word',179,'新增敏感词：网络水军','2026-06-04 13:36:04.313'),(215,1,'create_sensitive_word','sensitive_word',180,'新增敏感词：有偿删帖','2026-06-04 13:36:04.319'),(216,1,'create_sensitive_word','sensitive_word',181,'新增敏感词：网络推广','2026-06-04 13:36:04.324'),(217,1,'create_sensitive_word','sensitive_word',182,'新增敏感词：群发广告','2026-06-04 13:36:04.329'),(218,1,'create_sensitive_word','sensitive_word',183,'新增敏感词：诚招代理','2026-06-04 13:36:04.334'),(219,1,'create_sensitive_word','sensitive_word',184,'新增敏感词：工资日结','2026-06-04 13:36:04.340'),(220,1,'create_sensitive_word','sensitive_word',185,'新增敏感词：轻松兼职','2026-06-04 13:36:04.345'),(221,1,'create_sensitive_word','sensitive_word',186,'新增敏感词：淘宝刷单','2026-06-04 13:36:04.350'),(222,1,'create_sensitive_word','sensitive_word',187,'新增敏感词：刷信誉','2026-06-04 13:36:04.356'),(223,1,'create_sensitive_word','sensitive_word',188,'新增敏感词：引流推广','2026-06-04 13:36:04.362'),(224,1,'create_sensitive_word','sensitive_word',189,'新增敏感词：互粉','2026-06-04 13:36:04.368'),(225,1,'create_sensitive_word','sensitive_word',190,'新增敏感词：互关','2026-06-04 13:36:04.373'),(226,1,'create_sensitive_word','sensitive_word',191,'新增敏感词：刷单兼职','2026-06-04 13:36:04.379'),(227,1,'create_sensitive_word','sensitive_word',192,'新增敏感词：日结兼职','2026-06-04 13:36:04.385'),(228,1,'create_sensitive_word','sensitive_word',193,'新增敏感词：宝妈兼职','2026-06-04 13:36:04.392'),(229,1,'create_sensitive_word','sensitive_word',194,'新增敏感词：招打字员','2026-06-04 13:36:04.399'),(230,1,'create_sensitive_word','sensitive_word',195,'新增敏感词：缴纳入职费','2026-06-04 13:36:04.405'),(231,1,'create_sensitive_word','sensitive_word',196,'新增敏感词：先交钱','2026-06-04 13:36:04.412'),(232,1,'create_sensitive_word','sensitive_word',197,'新增敏感词：培训费','2026-06-04 13:36:04.417');
/*!40000 ALTER TABLE `operation_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(120) NOT NULL,
  `summary` varchar(300) DEFAULT NULL,
  `content` longtext NOT NULL,
  `cover_image` varchar(255) DEFAULT NULL,
  `tech_stacks` longtext,
  `highlights` longtext,
  `demo_url` varchar(255) DEFAULT NULL,
  `repo_url` varchar(255) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'draft',
  `is_featured` tinyint(1) NOT NULL DEFAULT '0',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `author_id` bigint unsigned NOT NULL,
  `published_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `role_label` varchar(120) DEFAULT NULL,
  `duration` varchar(120) DEFAULT NULL,
  `team_label` varchar(120) DEFAULT NULL,
  `process` longtext,
  `challenges` longtext,
  `solutions` longtext,
  `results` longtext,
  `reject_reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_projects_status` (`status`),
  KEY `idx_projects_author_id` (`author_id`),
  CONSTRAINT `fk_projects_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES (1,'Test Project Draft','Testing','Test project content','','[]','[]','','','published',0,0,1,'2026-06-01 12:38:25.711','2026-06-01 12:38:25.711','2026-06-01 12:38:25.711','','','','[]','[]','[]','[]',''),(2,'My Draft Project','Draft summary','Draft content','','[]','[]','','','draft',0,0,5,NULL,'2026-06-01 12:38:36.018','2026-06-01 12:38:36.018','','','','[]','[]','[]','[]',''),(3,'Updated','Updated','Updated','',NULL,NULL,'','','draft',0,0,1,NULL,'2026-06-01 12:38:50.114','2026-06-01 12:38:50.372','','','',NULL,NULL,NULL,NULL,''),(4,'QA������Ŀ','��Ŀ����','# ��Ŀ����\n\n����һ��������Ŀ','','[]','[]','https://demo.example.com','','published',0,0,1,'2026-06-01 14:09:49.161','2026-06-01 14:09:49.161','2026-06-01 14:09:49.161','','','','[]','[]','[]','[]',''),(5,'User Project Draft','','# User project','','[]','[]','','','draft',0,0,6,NULL,'2026-06-01 14:09:49.587','2026-06-01 14:09:49.587','','','','[]','[]','[]','[]',''),(6,'QA Test Project','Summary','# Content','','[]','[]','','','published',0,0,1,'2026-06-01 14:09:56.153','2026-06-01 14:09:56.153','2026-06-01 14:09:56.153','','','','[]','[]','[]','[]','');
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sensitive_words`
--

DROP TABLE IF EXISTS `sensitive_words`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sensitive_words` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `word` varchar(120) NOT NULL,
  `category` varchar(40) NOT NULL DEFAULT 'custom',
  `note` varchar(255) DEFAULT NULL,
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_sensitive_words_word` (`word`),
  KEY `idx_sensitive_words_category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=198 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sensitive_words`
--

LOCK TABLES `sensitive_words` WRITE;
/*!40000 ALTER TABLE `sensitive_words` DISABLE KEYS */;
INSERT INTO `sensitive_words` VALUES (1,'六四事件','political','历史敏感事件',1,'2026-06-04 13:34:40.043','2026-06-04 13:34:40.043'),(2,'天安门母亲','political','敏感组织',1,'2026-06-04 13:34:40.054','2026-06-04 13:34:40.054'),(3,'极端宗教','political','极端势力',1,'2026-06-04 13:34:40.060','2026-06-04 13:34:40.060'),(4,'宗教极端','political','极端势力',1,'2026-06-04 13:34:40.065','2026-06-04 13:34:40.065'),(5,'民族分裂','political','分裂主义',1,'2026-06-04 13:34:40.071','2026-06-04 13:34:40.071'),(6,'宗教迫害','political','宗教敏感',1,'2026-06-04 13:34:40.075','2026-06-04 13:34:40.075'),(7,'非法组织','political','违法组织',1,'2026-06-04 13:34:40.080','2026-06-04 13:34:40.080'),(8,'非法集会','political','群体事件',1,'2026-06-04 13:34:40.085','2026-06-04 13:34:40.085'),(9,'非法游行','political','群体事件',1,'2026-06-04 13:34:40.091','2026-06-04 13:34:40.091'),(10,'非法示威','political','群体事件',1,'2026-06-04 13:34:40.097','2026-06-04 13:34:40.097'),(11,'非法聚集','political','群体事件',1,'2026-06-04 13:34:40.102','2026-06-04 13:34:40.102'),(12,'静坐','political','抗议行为',1,'2026-06-04 13:34:40.107','2026-06-04 13:34:40.107'),(13,'拉横幅','political','抗议行为',1,'2026-06-04 13:34:40.112','2026-06-04 13:34:40.112'),(14,'越级上访','political','信访敏感',1,'2026-06-04 13:34:40.115','2026-06-04 13:34:40.115'),(15,'被消失','political','敏感表述',1,'2026-06-04 13:34:40.121','2026-06-04 13:34:40.121'),(16,'被自杀','political','敏感表述',1,'2026-06-04 13:34:40.126','2026-06-04 13:34:40.126'),(17,'恶意讨薪','political','维稳敏感',1,'2026-06-04 13:34:40.131','2026-06-04 13:34:40.131'),(18,'恐怖主义','violent','暴力恐怖',1,'2026-06-04 13:34:40.135','2026-06-04 13:34:40.135'),(19,'极端主义','violent','暴力恐怖',1,'2026-06-04 13:34:40.141','2026-06-04 13:34:40.141'),(20,'自焚','violent','自残行为',1,'2026-06-04 13:34:40.146','2026-06-04 13:34:40.146'),(21,'割腕','violent','自残行为',1,'2026-06-04 13:34:40.152','2026-06-04 13:34:40.152'),(22,'服毒','violent','自杀行为',1,'2026-06-04 13:34:40.158','2026-06-04 13:34:40.158'),(23,'安眠药','violent','自杀工具',1,'2026-06-04 13:34:40.163','2026-06-04 13:34:40.163'),(24,'百草枯','violent','剧毒农药',1,'2026-06-04 13:34:40.167','2026-06-04 13:34:40.167'),(25,'氰化物','violent','剧毒化学品',1,'2026-06-04 13:34:40.173','2026-06-04 13:34:40.173'),(26,'杀人灭口','violent','暴力犯罪',1,'2026-06-04 13:34:40.177','2026-06-04 13:34:40.177'),(27,'绑架','violent','暴力犯罪',1,'2026-06-04 13:34:40.184','2026-06-04 13:34:40.184'),(28,'分尸','violent','暴力犯罪',1,'2026-06-04 13:34:40.188','2026-06-04 13:34:40.188'),(29,'碎尸','violent','暴力犯罪',1,'2026-06-04 13:34:40.194','2026-06-04 13:34:40.194'),(30,'肢解','violent','暴力犯罪',1,'2026-06-04 13:34:40.198','2026-06-04 13:34:40.198'),(31,'斩首','violent','暴力恐怖视频',1,'2026-06-04 13:34:40.202','2026-06-04 13:34:40.202'),(32,'暴恐音视频','violent','暴力恐怖',1,'2026-06-04 13:34:40.208','2026-06-04 13:34:40.208'),(33,'血腥暴力','violent','暴力内容',1,'2026-06-04 13:34:40.213','2026-06-04 13:34:40.213'),(34,'迷药','contraband','违禁药物',1,'2026-06-04 13:34:40.218','2026-06-04 13:34:40.218'),(35,'迷魂药','contraband','违禁药物',1,'2026-06-04 13:34:40.228','2026-06-04 13:34:40.228'),(36,'催情药','contraband','违禁药物',1,'2026-06-04 13:34:40.235','2026-06-04 13:34:40.235'),(37,'听话水','contraband','违禁药物',1,'2026-06-04 13:34:40.241','2026-06-04 13:34:40.241'),(38,'神仙水','contraband','违禁药物',1,'2026-06-04 13:34:40.246','2026-06-04 13:34:40.246'),(39,'摇头丸','contraband','新型毒品',1,'2026-06-04 13:34:40.252','2026-06-04 13:34:40.252'),(40,'麻古','contraband','新型毒品',1,'2026-06-04 13:34:40.258','2026-06-04 13:34:40.258'),(41,'止咳水','contraband','滥用药物',1,'2026-06-04 13:34:40.264','2026-06-04 13:34:40.264'),(42,'蓝精灵','contraband','新型毒品',1,'2026-06-04 13:34:40.271','2026-06-04 13:34:40.271'),(43,'邮票毒品','contraband','新型毒品',1,'2026-06-04 13:34:40.276','2026-06-04 13:34:40.276'),(44,'电子烟毒品','contraband','新型毒品',1,'2026-06-04 13:34:40.282','2026-06-04 13:34:40.282'),(45,'制毒','contraband','制毒相关',1,'2026-06-04 13:34:40.287','2026-06-04 13:34:40.287'),(46,'贩毒','contraband','贩毒相关',1,'2026-06-04 13:34:40.295','2026-06-04 13:34:40.295'),(47,'仿真枪','contraband','违禁武器',1,'2026-06-04 13:34:40.300','2026-06-04 13:34:40.300'),(48,'气枪','contraband','违禁武器',1,'2026-06-04 13:34:40.304','2026-06-04 13:34:40.304'),(49,'猎枪','contraband','违禁武器',1,'2026-06-04 13:34:40.311','2026-06-04 13:34:40.311'),(50,'钢珠枪','contraband','违禁武器',1,'2026-06-04 13:34:40.316','2026-06-04 13:34:40.316'),(51,'电击棒','contraband','管制器具',1,'2026-06-04 13:34:40.322','2026-06-04 13:34:40.322'),(52,'弹簧刀','contraband','管制刀具',1,'2026-06-04 13:34:40.326','2026-06-04 13:34:40.326'),(53,'匕首','contraband','管制刀具',1,'2026-06-04 13:34:40.331','2026-06-04 13:34:40.331'),(54,'卖肾','contraband','非法器官交易',1,'2026-06-04 13:34:40.336','2026-06-04 13:34:40.336'),(55,'器官买卖','contraband','非法器官交易',1,'2026-06-04 13:34:40.341','2026-06-04 13:34:40.341'),(56,'代孕','contraband','非法生殖服务',1,'2026-06-04 13:34:40.345','2026-06-04 13:34:40.345'),(57,'卵子买卖','contraband','非法生殖服务',1,'2026-06-04 13:34:40.351','2026-06-04 13:34:40.351'),(58,'三级片','porn','淫秽色情',1,'2026-06-04 13:35:03.301','2026-06-04 13:35:03.301'),(59,'成人片','porn','淫秽色情',1,'2026-06-04 13:35:03.310','2026-06-04 13:35:03.310'),(60,'毛片','porn','淫秽色情',1,'2026-06-04 13:35:03.317','2026-06-04 13:35:03.317'),(61,'无码','porn','淫秽色情',1,'2026-06-04 13:35:03.322','2026-06-04 13:35:03.322'),(62,'成人电影','porn','淫秽色情',1,'2026-06-04 13:35:03.327','2026-06-04 13:35:03.327'),(63,'情色','porn','淫秽色情',1,'2026-06-04 13:35:03.332','2026-06-04 13:35:03.332'),(64,'裸模','porn','低俗色情',1,'2026-06-04 13:35:03.337','2026-06-04 13:35:03.337'),(65,'一夜情','porn','低俗色情',1,'2026-06-04 13:35:03.342','2026-06-04 13:35:03.342'),(66,'同城约炮','porn','低俗色情',1,'2026-06-04 13:35:03.347','2026-06-04 13:35:03.347'),(67,'淫秽','porn','淫秽色情',1,'2026-06-04 13:35:03.352','2026-06-04 13:35:03.352'),(68,'淫乱','porn','淫秽色情',1,'2026-06-04 13:35:03.357','2026-06-04 13:35:03.357'),(69,'淫荡','porn','淫秽色情',1,'2026-06-04 13:35:03.361','2026-06-04 13:35:03.361'),(70,'性交易','porn','违法交易',1,'2026-06-04 13:35:03.367','2026-06-04 13:35:03.367'),(71,'包养','porn','低俗色情',1,'2026-06-04 13:35:03.371','2026-06-04 13:35:03.371'),(72,'失足少女','porn','低俗暗示',1,'2026-06-04 13:35:03.375','2026-06-04 13:35:03.375'),(73,'陪睡','porn','色情服务',1,'2026-06-04 13:35:03.380','2026-06-04 13:35:03.380'),(74,'出台','porn','色情服务',1,'2026-06-04 13:35:03.384','2026-06-04 13:35:03.384'),(75,'坐台','porn','色情服务',1,'2026-06-04 13:35:03.389','2026-06-04 13:35:03.389'),(76,'裸体','porn','低俗内容',1,'2026-06-04 13:35:03.394','2026-06-04 13:35:03.394'),(77,'裸照','porn','低俗内容',1,'2026-06-04 13:35:03.399','2026-06-04 13:35:03.399'),(78,'偷拍','porn','低俗色情',1,'2026-06-04 13:35:03.404','2026-06-04 13:35:03.404'),(79,'走光','porn','低俗色情',1,'2026-06-04 13:35:03.410','2026-06-04 13:35:03.410'),(80,'艳照','porn','低俗色情',1,'2026-06-04 13:35:03.415','2026-06-04 13:35:03.415'),(81,'迷奸','porn','违法犯罪',1,'2026-06-04 13:35:03.420','2026-06-04 13:35:03.420'),(82,'兽交','porn','淫秽内容',1,'2026-06-04 13:35:03.426','2026-06-04 13:35:03.426'),(83,'乱伦','porn','淫秽内容',1,'2026-06-04 13:35:03.431','2026-06-04 13:35:03.431'),(84,'幼女','porn','儿童色情',1,'2026-06-04 13:35:03.435','2026-06-04 13:35:03.435'),(85,'少女色情','porn','儿童色情',1,'2026-06-04 13:35:03.441','2026-06-04 13:35:03.441'),(86,'萝莉','porn','儿童色情暗示',1,'2026-06-04 13:35:03.447','2026-06-04 13:35:03.447'),(87,'性奴','porn','淫秽内容',1,'2026-06-04 13:35:03.452','2026-06-04 13:35:03.452'),(88,'sm','porn','淫秽内容',1,'2026-06-04 13:35:03.457','2026-06-04 13:35:03.457'),(89,'nsfw','porn','色情标签',1,'2026-06-04 13:35:03.462','2026-06-04 13:35:03.462'),(90,'onlyfans','porn','成人平台',1,'2026-06-04 13:35:03.468','2026-06-04 13:35:03.468'),(91,'时时彩','gambling','非法彩票',1,'2026-06-04 13:35:03.472','2026-06-04 13:35:03.472'),(92,'六合彩','gambling','非法彩票',1,'2026-06-04 13:35:03.477','2026-06-04 13:35:03.477'),(93,'老虎机','gambling','赌博设备',1,'2026-06-04 13:35:03.482','2026-06-04 13:35:03.482'),(94,'百家乐','gambling','赌博形式',1,'2026-06-04 13:35:03.486','2026-06-04 13:35:03.486'),(95,'德州扑克','gambling','赌博形式',1,'2026-06-04 13:35:03.492','2026-06-04 13:35:03.492'),(96,'牌九','gambling','赌博形式',1,'2026-06-04 13:35:03.497','2026-06-04 13:35:03.497'),(97,'梭哈','gambling','赌博形式',1,'2026-06-04 13:35:03.501','2026-06-04 13:35:03.501'),(98,'赌球','gambling','赌博形式',1,'2026-06-04 13:35:03.506','2026-06-04 13:35:03.506'),(99,'赌马','gambling','赌博形式',1,'2026-06-04 13:35:03.511','2026-06-04 13:35:03.511'),(100,'赌船','gambling','赌博场所',1,'2026-06-04 13:35:03.516','2026-06-04 13:35:03.516'),(101,'网络赌场','gambling','赌博平台',1,'2026-06-04 13:35:03.522','2026-06-04 13:35:03.522'),(102,'棋牌赌博','gambling','赌博形式',1,'2026-06-04 13:35:03.527','2026-06-04 13:35:03.527'),(103,'捕鱼游戏','gambling','赌博游戏',1,'2026-06-04 13:35:03.535','2026-06-04 13:35:03.535'),(104,'电子游艺','gambling','赌博游戏',1,'2026-06-04 13:35:03.540','2026-06-04 13:35:03.540'),(105,'体育投注','gambling','赌博形式',1,'2026-06-04 13:35:03.546','2026-06-04 13:35:03.546'),(106,'电子赌博','gambling','赌博形式',1,'2026-06-04 13:35:03.552','2026-06-04 13:35:03.552'),(107,'网络赌博','gambling','赌博平台',1,'2026-06-04 13:35:03.559','2026-06-04 13:35:03.559'),(108,'投注平台','gambling','赌博网站',1,'2026-06-04 13:35:03.565','2026-06-04 13:35:03.565'),(109,'现金棋牌','gambling','赌博平台',1,'2026-06-04 13:35:03.570','2026-06-04 13:35:03.570'),(110,'信誉平台','gambling','赌博宣传',1,'2026-06-04 13:35:03.576','2026-06-04 13:35:03.576'),(111,'真人娱乐','gambling','赌博宣传',1,'2026-06-04 13:35:03.580','2026-06-04 13:35:03.580'),(112,'澳门赌场','gambling','赌博场所',1,'2026-06-04 13:35:03.584','2026-06-04 13:35:03.584'),(113,'葡京','gambling','赌博场所',1,'2026-06-04 13:35:03.593','2026-06-04 13:35:03.593'),(114,'网络传销','scam','诈骗传销',1,'2026-06-04 13:35:03.598','2026-06-04 13:35:03.598'),(115,'资金盘','scam','金融诈骗',1,'2026-06-04 13:35:03.603','2026-06-04 13:35:03.603'),(116,'杀猪盘','scam','恋爱诈骗',1,'2026-06-04 13:35:03.607','2026-06-04 13:35:03.607'),(117,'庞氏骗局','scam','金融诈骗',1,'2026-06-04 13:35:03.612','2026-06-04 13:35:03.612'),(118,'拉人头','scam','传销模式',1,'2026-06-04 13:35:03.620','2026-06-04 13:35:03.620'),(119,'零风险高收益','scam','诈骗宣传',1,'2026-06-04 13:35:03.624','2026-06-04 13:35:03.624'),(120,'稳赚不赔','scam','诈骗宣传',1,'2026-06-04 13:35:03.630','2026-06-04 13:35:03.630'),(121,'日赚','scam','诈骗广告',1,'2026-06-04 13:35:03.635','2026-06-04 13:35:03.635'),(122,'月入过万','scam','诈骗广告',1,'2026-06-04 13:35:03.641','2026-06-04 13:35:03.641'),(123,'挂机赚钱','scam','诈骗广告',1,'2026-06-04 13:35:03.646','2026-06-04 13:35:03.646'),(124,'点赞赚钱','scam','刷单诈骗',1,'2026-06-04 13:35:03.651','2026-06-04 13:35:03.651'),(125,'虚拟币诈骗','scam','金融诈骗',1,'2026-06-04 13:35:03.657','2026-06-04 13:35:03.657'),(126,'区块链骗局','scam','金融诈骗',1,'2026-06-04 13:35:03.661','2026-06-04 13:35:03.661'),(127,'投资返利','scam','诈骗宣传',1,'2026-06-04 13:35:03.668','2026-06-04 13:35:03.668'),(128,'推荐股票','scam','投资诈骗',1,'2026-06-04 13:35:03.673','2026-06-04 13:35:03.673'),(129,'内幕消息','scam','投资诈骗',1,'2026-06-04 13:35:03.680','2026-06-04 13:35:03.680'),(130,'割韭菜','scam','投资诈骗',1,'2026-06-04 13:35:03.685','2026-06-04 13:35:03.685'),(131,'薅羊毛','scam','刷单诈骗',1,'2026-06-04 13:35:03.690','2026-06-04 13:35:03.690'),(132,'接码平台','scam','黑产工具',1,'2026-06-04 13:35:03.696','2026-06-04 13:35:03.696'),(133,'猫池','scam','黑产设备',1,'2026-06-04 13:35:03.701','2026-06-04 13:35:03.701'),(134,'伪基站','scam','黑产设备',1,'2026-06-04 13:35:03.706','2026-06-04 13:35:03.706'),(135,'短信轰炸','scam','黑产攻击',1,'2026-06-04 13:35:03.712','2026-06-04 13:35:03.712'),(136,'电话轰炸','scam','黑产攻击',1,'2026-06-04 13:35:03.717','2026-06-04 13:35:03.717'),(137,'种族歧视','hate','歧视内容',1,'2026-06-04 13:36:04.069','2026-06-04 13:36:04.069'),(138,'地域黑','hate','地域歧视',1,'2026-06-04 13:36:04.080','2026-06-04 13:36:04.080'),(139,'地域歧视','hate','地域歧视',1,'2026-06-04 13:36:04.087','2026-06-04 13:36:04.087'),(140,'歧视女性','hate','性别歧视',1,'2026-06-04 13:36:04.093','2026-06-04 13:36:04.093'),(141,'歧视残疾人','hate','残障歧视',1,'2026-06-04 13:36:04.099','2026-06-04 13:36:04.099'),(142,'人身攻击','hate','攻击辱骂',1,'2026-06-04 13:36:04.105','2026-06-04 13:36:04.105'),(143,'网络暴力','hate','网络暴力',1,'2026-06-04 13:36:04.110','2026-06-04 13:36:04.110'),(144,'人肉搜索','hate','隐私侵犯',1,'2026-06-04 13:36:04.116','2026-06-04 13:36:04.116'),(145,'开盒','hate','隐私侵犯',1,'2026-06-04 13:36:04.122','2026-06-04 13:36:04.122'),(146,'死亡威胁','hate','暴力威胁',1,'2026-06-04 13:36:04.128','2026-06-04 13:36:04.128'),(147,'去死','hate','恶意诅咒',1,'2026-06-04 13:36:04.134','2026-06-04 13:36:04.134'),(148,'全家死','hate','恶意诅咒',1,'2026-06-04 13:36:04.139','2026-06-04 13:36:04.139'),(149,'畜生','hate','侮辱谩骂',1,'2026-06-04 13:36:04.145','2026-06-04 13:36:04.145'),(150,'娘炮','hate','性别歧视',1,'2026-06-04 13:36:04.151','2026-06-04 13:36:04.151'),(151,'脱库','illegal','数据泄露',1,'2026-06-04 13:36:04.156','2026-06-04 13:36:04.156'),(152,'撞库','illegal','数据泄露',1,'2026-06-04 13:36:04.162','2026-06-04 13:36:04.162'),(153,'删库跑路','illegal','破坏行为',1,'2026-06-04 13:36:04.167','2026-06-04 13:36:04.167'),(154,'webshell','illegal','黑客工具',1,'2026-06-04 13:36:04.172','2026-06-04 13:36:04.172'),(155,'卷款潜逃','illegal','经济犯罪',1,'2026-06-04 13:36:04.178','2026-06-04 13:36:04.178'),(156,'非法经营','illegal','违法经营',1,'2026-06-04 13:36:04.183','2026-06-04 13:36:04.183'),(157,'无证经营','illegal','违法经营',1,'2026-06-04 13:36:04.189','2026-06-04 13:36:04.189'),(158,'偷税漏税','illegal','税务违法',1,'2026-06-04 13:36:04.195','2026-06-04 13:36:04.195'),(159,'走私','illegal','走私犯罪',1,'2026-06-04 13:36:04.200','2026-06-04 13:36:04.200'),(160,'地下钱庄','illegal','非法金融',1,'2026-06-04 13:36:04.204','2026-06-04 13:36:04.204'),(161,'高利贷','illegal','非法借贷',1,'2026-06-04 13:36:04.209','2026-06-04 13:36:04.209'),(162,'套路贷','illegal','非法借贷',1,'2026-06-04 13:36:04.215','2026-06-04 13:36:04.215'),(163,'校园贷','illegal','非法借贷',1,'2026-06-04 13:36:04.221','2026-06-04 13:36:04.221'),(164,'裸贷','illegal','非法借贷',1,'2026-06-04 13:36:04.226','2026-06-04 13:36:04.226'),(165,'暴力催收','illegal','暴力催收',1,'2026-06-04 13:36:04.232','2026-06-04 13:36:04.232'),(166,'盗版资源','illegal','版权侵权',1,'2026-06-04 13:36:04.237','2026-06-04 13:36:04.237'),(167,'注册机','illegal','软件破解',1,'2026-06-04 13:36:04.243','2026-06-04 13:36:04.243'),(168,'游戏外挂','illegal','游戏作弊',1,'2026-06-04 13:36:04.248','2026-06-04 13:36:04.248'),(169,'私服','illegal','私服运营',1,'2026-06-04 13:36:04.253','2026-06-04 13:36:04.253'),(170,'电信诈骗','illegal','电信诈骗',1,'2026-06-04 13:36:04.260','2026-06-04 13:36:04.260'),(171,'帮信罪','illegal','帮助信息网络犯罪',1,'2026-06-04 13:36:04.266','2026-06-04 13:36:04.266'),(172,'跨境赌博','illegal','跨境犯罪',1,'2026-06-04 13:36:04.272','2026-06-04 13:36:04.272'),(173,'代写论文','spam','学术不端',1,'2026-06-04 13:36:04.277','2026-06-04 13:36:04.277'),(174,'代发论文','spam','学术不端',1,'2026-06-04 13:36:04.283','2026-06-04 13:36:04.283'),(175,'论文代写','spam','学术不端',1,'2026-06-04 13:36:04.288','2026-06-04 13:36:04.288'),(176,'刷粉','spam','刷量服务',1,'2026-06-04 13:36:04.293','2026-06-04 13:36:04.293'),(177,'刷赞','spam','刷量服务',1,'2026-06-04 13:36:04.297','2026-06-04 13:36:04.297'),(178,'刷屏','spam','恶意刷屏',1,'2026-06-04 13:36:04.305','2026-06-04 13:36:04.305'),(179,'网络水军','spam','网络水军',1,'2026-06-04 13:36:04.310','2026-06-04 13:36:04.310'),(180,'有偿删帖','spam','黑灰产',1,'2026-06-04 13:36:04.316','2026-06-04 13:36:04.316'),(181,'网络推广','spam','垃圾广告',1,'2026-06-04 13:36:04.322','2026-06-04 13:36:04.322'),(182,'群发广告','spam','垃圾广告',1,'2026-06-04 13:36:04.327','2026-06-04 13:36:04.327'),(183,'诚招代理','spam','微商广告',1,'2026-06-04 13:36:04.331','2026-06-04 13:36:04.331'),(184,'工资日结','spam','兼职诈骗广告',1,'2026-06-04 13:36:04.337','2026-06-04 13:36:04.337'),(185,'轻松兼职','spam','兼职诈骗广告',1,'2026-06-04 13:36:04.342','2026-06-04 13:36:04.342'),(186,'淘宝刷单','spam','刷单广告',1,'2026-06-04 13:36:04.348','2026-06-04 13:36:04.348'),(187,'刷信誉','spam','刷单广告',1,'2026-06-04 13:36:04.354','2026-06-04 13:36:04.354'),(188,'引流推广','spam','流量引导',1,'2026-06-04 13:36:04.359','2026-06-04 13:36:04.359'),(189,'互粉','spam','互粉广告',1,'2026-06-04 13:36:04.365','2026-06-04 13:36:04.365'),(190,'互关','spam','互关广告',1,'2026-06-04 13:36:04.371','2026-06-04 13:36:04.371'),(191,'刷单兼职','spam','刷单诈骗',1,'2026-06-04 13:36:04.377','2026-06-04 13:36:04.377'),(192,'日结兼职','spam','兼职诈骗',1,'2026-06-04 13:36:04.381','2026-06-04 13:36:04.381'),(193,'宝妈兼职','spam','兼职诈骗',1,'2026-06-04 13:36:04.388','2026-06-04 13:36:04.388'),(194,'招打字员','spam','兼职诈骗',1,'2026-06-04 13:36:04.396','2026-06-04 13:36:04.396'),(195,'缴纳入职费','spam','招聘诈骗',1,'2026-06-04 13:36:04.402','2026-06-04 13:36:04.402'),(196,'先交钱','spam','招聘诈骗',1,'2026-06-04 13:36:04.409','2026-06-04 13:36:04.409'),(197,'培训费','spam','招聘诈骗',1,'2026-06-04 13:36:04.415','2026-06-04 13:36:04.415');
/*!40000 ALTER TABLE `sensitive_words` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `system_settings`
--

DROP TABLE IF EXISTS `system_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `system_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(80) NOT NULL,
  `value` varchar(255) NOT NULL,
  `note` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_system_settings_key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `system_settings`
--

LOCK TABLES `system_settings` WRITE;
/*!40000 ALTER TABLE `system_settings` DISABLE KEYS */;
INSERT INTO `system_settings` VALUES (1,'moderation_ban_threshold','3','24 小时内触发多少次违禁词后自动封禁普通用户','2026-04-15 17:27:19.541','2026-06-01 14:10:41.757');
/*!40000 ALTER TABLE `system_settings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tags` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tags_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tags`
--

LOCK TABLES `tags` WRITE;
/*!40000 ALTER TABLE `tags` DISABLE KEYS */;
INSERT INTO `tags` VALUES (1,'dwadawd','2026-04-02 13:38:36.886','2026-04-02 13:38:36.886'),(2,'vibe coding','2026-04-15 19:06:08.999','2026-04-15 19:06:08.999'),(3,'qa-test-tag','2026-06-01 14:09:50.024','2026-06-01 14:09:50.024'),(4,'qa-tag','2026-06-01 14:10:24.002','2026-06-01 14:10:24.002');
/*!40000 ALTER TABLE `tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` varchar(16) NOT NULL DEFAULT 'user',
  `bio` varchar(2000) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `status` varchar(16) NOT NULL DEFAULT 'active',
  `avatar` varchar(255) DEFAULT NULL,
  `headline` varchar(120) DEFAULT NULL,
  `current_role` varchar(120) DEFAULT NULL,
  `years_label` varchar(80) DEFAULT NULL,
  `motto` varchar(255) DEFAULT NULL,
  `location` varchar(120) DEFAULT NULL,
  `email` varchar(120) DEFAULT NULL,
  `resume_url` varchar(255) DEFAULT NULL,
  `website_url` varchar(255) DEFAULT NULL,
  `github_url` varchar(255) DEFAULT NULL,
  `gitee_url` varchar(255) DEFAULT NULL,
  `juejin_url` varchar(255) DEFAULT NULL,
  `csdn_url` varchar(255) DEFAULT NULL,
  `skills` longtext,
  `focus_areas` longtext,
  `ban_reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin','$2a$10$Wgofj82PL6CgSpgkcwmekuYL9NN4A32QZIXr9Jh0Nf0Wp/uEl5ffy','admin','���Ǹ��˼���','2026-04-01 21:16:22.810','2026-06-01 14:06:30.745','active','','ȫջ����ʦ','','','','����','','','','','','','','[\"Go\",\"Vue\",\"Docker\"]','[\"���˿���\",\"��ԭ��\"]',NULL),(2,'gtfh12','$2a$10$poR0oUyjw2i9JveFCMDGsusTGOcpGWhFB8GBORlC1oguYr/DfCi6q','user','这个站点会记录我的项目实践、技术笔记和成长轨迹，让学习与输出都变成能被看见的成果。','2026-04-01 21:27:40.580','2026-06-01 14:13:05.009','active','','通过项目、笔记和公开输出持续成长的开发者。','成长型开发者','在项目与写作中持续进步','认真学习，清楚表达，持续构建。','东莞','','','','','','','','[\"Vue 3\",\"Go\",\"笔记写作\",\"问题排查\"]','[\"学习归档\",\"项目作品集\",\"公开写作\"]',NULL),(3,'hhh','$2a$10$sfeyJY/Lw0rzrYimUaOTa.7A.QmW/3JHHbr1SM2ObSzigHkM7PcJa','user','','2026-04-01 22:52:05.005','2026-04-01 22:52:05.005','active',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(4,'凌杰源','$2a$10$PE/c2URdltiPeaHx8rIwuudjuwNnS//axk4NBatJP1OnsFuzje6Z.','admin','','2026-04-15 19:04:10.485','2026-04-15 19:41:56.814','active','','','','','','','','','','','','','',NULL,NULL,''),(5,'testuser1','$2a$10$eqTkPrm9Qsb9lebzpavNpeYBVnPGOWGTiJefEXNIgk0kg0vAL58HG','user','','2026-06-01 12:38:35.795','2026-06-01 12:38:35.795','active','','','','','','','','','','','','','',NULL,NULL,''),(6,'qauser','$2a$10$czoJB.HUbIvp5Q2egtNOjOaWc1BrZ3amtXxe/oyampSMgvHIr8N/6','user','','2026-06-01 14:03:20.492','2026-06-01 14:10:42.033','active','','','','','','','','','','','','','',NULL,NULL,''),(7,'pwtest4','$2a$10$WxkrqIhpkQL8H7UOflbsJu6QJH0FBSjhvwd4Hir3RUIGgyWQwQAPi','user','','2026-06-02 00:04:37.104','2026-06-02 00:04:37.104','active','','','','','','','','','','','','','',NULL,NULL,''),(8,'admin1','$2a$10$zvX28UupvOl9aTAVn4qsoeIuXkIM3TGYIQK/w3S4Snj6K2aCIKnY6','user','','2026-06-04 18:48:04.391','2026-06-04 18:48:04.391','active','','','','','','','','','','','','','',NULL,NULL,'');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-06 20:56:06
