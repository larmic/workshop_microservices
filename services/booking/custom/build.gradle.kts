plugins {
    kotlin("jvm") version "2.4.0"
    kotlin("plugin.serialization") version "2.4.0"
    application
    id("com.gradleup.shadow") version "9.5.1"
}

group = "workshop.booking"
version = "1.0.0"

repositories {
    mavenCentral()
}

val ktorVersion = "3.5.1"

dependencies {
    implementation("io.ktor:ktor-server-core:$ktorVersion")
    implementation("io.ktor:ktor-server-cio:$ktorVersion")
    implementation("io.ktor:ktor-server-content-negotiation:$ktorVersion")
    implementation("io.ktor:ktor-serialization-kotlinx-json:$ktorVersion")
    implementation("io.ktor:ktor-client-core:$ktorVersion")
    implementation("io.ktor:ktor-client-cio:$ktorVersion")
    implementation("ch.qos.logback:logback-classic:1.5.37")
}

application {
    mainClass.set("workshop.booking.ApplicationKt")
}

kotlin {
    jvmToolchain(25)
}

tasks.shadowJar {
    archiveBaseName.set("booking-custom")
    archiveClassifier.set("all")
    archiveVersion.set("")
    mergeServiceFiles()
}
